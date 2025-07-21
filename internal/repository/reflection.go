// Updated internal/repository/reflection.go
package repository

import (
    "database/sql"
    "strings"
    "reflection-app/internal/models"
    "log"
    "time"
    "fmt"
)

type ReflectionRepository struct {
    db *sql.DB
}

func NewReflectionRepository(db *sql.DB) *ReflectionRepository {
    return &ReflectionRepository{db: db}
}

func (r *ReflectionRepository) Create(reflection *models.Reflection) error {
    // Convert the date to a date-only string to avoid timezone conversion issues
    dateString := reflection.Date.Time().Format("2006-01-02")
    log.Printf("Create: Saving reflection with date string: %s (original: %s)", 
        dateString, reflection.Date.Time().Format("2006-01-02 15:04:05 MST"))
    
    // Set default visibility if not provided
    if reflection.Visibility == "" {
        if reflection.IsPrivate {
            reflection.Visibility = "private"
        } else {
            reflection.Visibility = "public"
        }
    }
    
    query := `
        INSERT INTO reflections (author_id, category_id, sub_category_id, date, reflection_text, reflection_detail, tags, is_private, visibility)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at`
    
    // Convert tags to PostgreSQL array format
    var tagsArray interface{}
    if len(reflection.Tags) == 0 {
        tagsArray = "{}"
    } else {
        // Convert to PostgreSQL array format: {"item1","item2","item3"}
        result := "{"
        for i, tag := range reflection.Tags {
            if i > 0 {
                result += ","
            }
            // Escape quotes and backslashes
            escaped := strings.ReplaceAll(tag, `\`, `\\`)
            escaped = strings.ReplaceAll(escaped, `"`, `\"`)
            result += `"` + escaped + `"`
        }
        result += "}"
        tagsArray = result
    }
    
    return r.db.QueryRow(
        query,
        reflection.AuthorID,
        reflection.CategoryID,
        reflection.SubCategoryID,
        dateString,
        reflection.ReflectionText,
        reflection.ReflectionDetail,
        tagsArray,
        reflection.IsPrivate,
        reflection.Visibility,
    ).Scan(&reflection.ID, &reflection.CreatedAt, &reflection.UpdatedAt)
}

func (r *ReflectionRepository) GetByUserID(userID int) ([]models.Reflection, error) {
    // Try the complex query with tracking data first
    query := `
        SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
               r.reflection_text, r.reflection_detail, r.tags, r.is_private, r.visibility,
               r.created_at, r.updated_at,
               c.name as category_name, c.description as category_desc,
               sc.name as subcategory_name, sc.description as subcategory_desc,
               COALESCE(rt_count.total_count, 0) as reflection_count,
               CASE WHEN rt_today.reflected_date IS NOT NULL THEN true ELSE false END as reflected_today
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
        LEFT JOIN (
            SELECT reflection_id, COUNT(*) as total_count
            FROM reflection_tracking
            WHERE user_id = $1
            GROUP BY reflection_id
        ) rt_count ON r.id = rt_count.reflection_id
        LEFT JOIN (
            SELECT reflection_id, reflected_date
            FROM reflection_tracking
            WHERE user_id = $1 AND reflected_date = $2
        ) rt_today ON r.id = rt_today.reflection_id
        WHERE r.author_id = $1 
        ORDER BY r.date DESC, r.created_at DESC`
    
    // Use current date in Go instead of PostgreSQL CURRENT_DATE to avoid timezone issues
    currentDate := time.Now().Format("2006-01-02")
    log.Printf("GetByUserID: Using current date: %s", currentDate)
    rows, err := r.db.Query(query, userID, currentDate)
    if err != nil {
        log.Printf("Error in GetByUserID query with tracking: %v", err)
        log.Printf("Falling back to basic query without tracking data...")
        
        // Fallback to basic query without tracking data
        basicQuery := `
            SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
                   r.reflection_text, r.reflection_detail, r.tags, r.is_private, r.visibility,
                   r.created_at, r.updated_at,
                   c.name as category_name, c.description as category_desc,
                   sc.name as subcategory_name, sc.description as subcategory_desc
            FROM reflections r
            JOIN categories c ON r.category_id = c.id
            LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
            WHERE r.author_id = $1 
            ORDER BY r.date DESC, r.created_at DESC`
        
        rows, err = r.db.Query(basicQuery, userID)
        if err != nil {
            log.Printf("Error in GetByUserID basic query: %v", err)
            return nil, err
        }
    }
    defer rows.Close()
    
    var reflections []models.Reflection
    for rows.Next() {
        var reflection models.Reflection
        var categoryName, categoryDesc string
        var subcategoryName, subcategoryDesc sql.NullString
        var reflectionCount int
        var reflectedToday bool
        
        // Try to scan with tracking data first
        err := rows.Scan(
            &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
            &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate, &reflection.Visibility,
            &reflection.CreatedAt, &reflection.UpdatedAt,
            &categoryName, &categoryDesc,
            &subcategoryName, &subcategoryDesc,
            &reflectionCount, &reflectedToday,
        )
        if err != nil {
            // If that fails, try without tracking data
            log.Printf("Error scanning with tracking data, trying without: %v", err)
            err = rows.Scan(
                &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
                &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate, &reflection.Visibility,
                &reflection.CreatedAt, &reflection.UpdatedAt,
                &categoryName, &categoryDesc,
                &subcategoryName, &subcategoryDesc,
            )
            if err != nil {
                log.Printf("Error scanning reflection row: %v", err)
                return nil, err
            }
            // Set default tracking data
            reflectionCount = 0
            reflectedToday = false
        }
        
        // Set tracking data
        reflection.ReflectionCount = reflectionCount
        reflection.ReflectedToday = reflectedToday
        
        log.Printf("GetByUserID: Retrieved reflection %d with date: %s (Raw from DB: %s)", 
            reflection.ID,
            reflection.Date.Time().Format("2006-01-02"),
            reflection.Date.Time().Format("2006-01-02 15:04:05 MST"))
        
        // Populate category
        reflection.Category = &models.Category{
            ID:          reflection.CategoryID,
            Name:        categoryName,
            Description: categoryDesc,
        }
        
        // Populate subcategory if exists
        if subcategoryName.Valid {
            reflection.SubCategory = &models.SubCategory{
                ID:          *reflection.SubCategoryID,
                CategoryID:  reflection.CategoryID,
                Name:        subcategoryName.String,
                Description: subcategoryDesc.String,
            }
        }
        
        // Fetch action items for this reflection
        actionsQuery := `
            SELECT id, reflection_id, action, priority, status, due_date, created_at, updated_at
            FROM actions 
            WHERE reflection_id = $1 
            ORDER BY created_at ASC`
        
        actionRows, err := r.db.Query(actionsQuery, reflection.ID)
        if err != nil {
            log.Printf("Error fetching actions for reflection %d: %v", reflection.ID, err)
            return nil, err
        }
        
        var actions []models.ActionItem
        for actionRows.Next() {
            var action models.ActionItem
            err := actionRows.Scan(
                &action.ID,
                &action.ReflectionID,
                &action.Action,
                &action.Priority,
                &action.Status,
                &action.DueDate,
                &action.CreatedAt,
                &action.UpdatedAt,
            )
            if err != nil {
                actionRows.Close()
                log.Printf("Error scanning action row: %v", err)
                return nil, err
            }
            actions = append(actions, action)
        }
        actionRows.Close()
        
        reflection.Actions = actions
        reflections = append(reflections, reflection)
    }
    
    return reflections, nil
}

func (r *ReflectionRepository) GetByID(id int) (*models.Reflection, error) {
    // First get the basic reflection data to get the author_id
    basicQuery := `
        SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
               r.reflection_text, r.reflection_detail, r.tags, r.is_private, 
               r.created_at, r.updated_at,
               c.name as category_name, c.description as category_desc,
               sc.name as subcategory_name, sc.description as subcategory_desc
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
        WHERE r.id = $1`
    
    reflection := &models.Reflection{}
    var categoryName, categoryDesc string
    var subcategoryName, subcategoryDesc sql.NullString
    
    err := r.db.QueryRow(basicQuery, id).Scan(
        &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
        &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate,
        &reflection.CreatedAt, &reflection.UpdatedAt,
        &categoryName, &categoryDesc,
        &subcategoryName, &subcategoryDesc,
    )
    
    if err != nil {
        return nil, err
    }
    
    // Get tracking data separately
    trackingQuery := `
        SELECT 
            COALESCE(rt_count.total_count, 0) as reflection_count,
            CASE WHEN rt_today.reflected_date IS NOT NULL THEN true ELSE false END as reflected_today
        FROM reflections r
        LEFT JOIN (
            SELECT reflection_id, COUNT(*) as total_count
            FROM reflection_tracking
            WHERE user_id = $1
            GROUP BY reflection_id
        ) rt_count ON r.id = rt_count.reflection_id
        LEFT JOIN (
            SELECT reflection_id, reflected_date
            FROM reflection_tracking
            WHERE user_id = $1 AND reflected_date = CURRENT_DATE
        ) rt_today ON r.id = rt_today.reflection_id
        WHERE r.id = $2`
    
    var reflectionCount int
    var reflectedToday bool
    
    err = r.db.QueryRow(trackingQuery, reflection.AuthorID, id).Scan(&reflectionCount, &reflectedToday)
    if err != nil {
        // If tracking query fails, set defaults
        reflectionCount = 0
        reflectedToday = false
    }
    
    // Set tracking data
    reflection.ReflectionCount = reflectionCount
    reflection.ReflectedToday = reflectedToday
    
    // Populate category
    reflection.Category = &models.Category{
        ID:          reflection.CategoryID,
        Name:        categoryName,
        Description: categoryDesc,
    }
    
    // Populate subcategory if exists
    if subcategoryName.Valid {
        reflection.SubCategory = &models.SubCategory{
            ID:          *reflection.SubCategoryID,
            CategoryID:  reflection.CategoryID,
            Name:        subcategoryName.String,
            Description: subcategoryDesc.String,
        }
    }
    
    // Fetch action items for this reflection
    actionsQuery := `
        SELECT id, reflection_id, action, priority, status, due_date, created_at, updated_at
        FROM actions 
        WHERE reflection_id = $1 
        ORDER BY created_at ASC`
    
    actionRows, err := r.db.Query(actionsQuery, reflection.ID)
    if err != nil {
        return nil, err
    }
    defer actionRows.Close()
    
    var actions []models.ActionItem
    for actionRows.Next() {
        var action models.ActionItem
        err := actionRows.Scan(
            &action.ID,
            &action.ReflectionID,
            &action.Action,
            &action.Priority,
            &action.Status,
            &action.DueDate,
            &action.CreatedAt,
            &action.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        actions = append(actions, action)
    }
    
    reflection.Actions = actions
    return reflection, nil
}

func (r *ReflectionRepository) Update(reflection *models.Reflection) error {
    // Convert the date to a date-only string to avoid timezone conversion issues
    dateString := reflection.Date.Time().Format("2006-01-02")
    log.Printf("Update: Saving reflection with date string: %s (original: %s)", 
        dateString, reflection.Date.Time().Format("2006-01-02 15:04:05 MST"))
    
    // Set visibility based on is_private if not provided
    if reflection.Visibility == "" {
        if reflection.IsPrivate {
            reflection.Visibility = "private"
        } else {
            reflection.Visibility = "public"
        }
    }
    
    query := `
        UPDATE reflections 
        SET category_id = $1, sub_category_id = $2, date = $3, reflection_text = $4, 
            reflection_detail = $5, tags = $6, is_private = $7, visibility = $8, updated_at = CURRENT_TIMESTAMP
        WHERE id = $9`
    
    // Convert tags to PostgreSQL array format
    var tagsArray interface{}
    if len(reflection.Tags) == 0 {
        tagsArray = "{}"
    } else {
        // Convert to PostgreSQL array format: {"item1","item2","item3"}
        result := "{"
        for i, tag := range reflection.Tags {
            if i > 0 {
                result += ","
            }
            // Escape quotes and backslashes
            escaped := strings.ReplaceAll(tag, `\`, `\\`)
            escaped = strings.ReplaceAll(escaped, `"`, `\"`)
            result += `"` + escaped + `"`
        }
        result += "}"
        tagsArray = result
    }
    
    _, err := r.db.Exec(
        query,
        reflection.CategoryID,
        reflection.SubCategoryID,
        dateString,
        reflection.ReflectionText,
        reflection.ReflectionDetail,
        tagsArray,
        reflection.IsPrivate,
        reflection.Visibility,
        reflection.ID,
    )
    
    return err
}

func (r *ReflectionRepository) Delete(id int) error {
    query := `DELETE FROM reflections WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}

func (r *ReflectionRepository) GetByCategory(userID, categoryID int) ([]models.Reflection, error) {
    query := `
        SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
               r.reflection_text, r.reflection_detail, r.tags, r.is_private, 
               r.created_at, r.updated_at,
               c.name as category_name, c.description as category_desc,
               sc.name as subcategory_name, sc.description as subcategory_desc
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
        WHERE r.author_id = $1 AND r.category_id = $2
        ORDER BY r.date DESC, r.created_at DESC`
    
    rows, err := r.db.Query(query, userID, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var reflections []models.Reflection
    for rows.Next() {
        var reflection models.Reflection
        var categoryName, categoryDesc string
        var subcategoryName, subcategoryDesc sql.NullString
        
        err := rows.Scan(
            &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
            &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate,
            &reflection.CreatedAt, &reflection.UpdatedAt,
            &categoryName, &categoryDesc,
            &subcategoryName, &subcategoryDesc,
        )
        if err != nil {
            return nil, err
        }
        
        // Populate category
        reflection.Category = &models.Category{
            ID:          reflection.CategoryID,
            Name:        categoryName,
            Description: categoryDesc,
        }
        
        // Populate subcategory if exists
        if subcategoryName.Valid {
            reflection.SubCategory = &models.SubCategory{
                ID:          *reflection.SubCategoryID,
                CategoryID:  reflection.CategoryID,
                Name:        subcategoryName.String,
                Description: subcategoryDesc.String,
            }
        }
        
        reflections = append(reflections, reflection)
    }
    
    return reflections, nil
}

func (r *ReflectionRepository) TrackReflection(reflectionID int) error {
    // Get the reflection to get the author_id without tracking data to avoid circular reference
    query := `
        SELECT author_id
        FROM reflections
        WHERE id = $1`
    
    var authorID int
    err := r.db.QueryRow(query, reflectionID).Scan(&authorID)
    if err != nil {
        return err
    }
    
    // Use current date in Go instead of PostgreSQL CURRENT_DATE to avoid timezone issues
    currentDate := time.Now().Format("2006-01-02")
    log.Printf("Tracking reflection %d for user %d with date: %s", reflectionID, authorID, currentDate)
    
    // Check if reflection is already tracked for today
    checkQuery := `
        SELECT id FROM reflection_tracking 
        WHERE reflection_id = $1 AND user_id = $2 AND reflected_date = $3`
    
    log.Printf("TrackReflection: Checking for existing tracking with reflection_id=%d, user_id=%d, date=%s", reflectionID, authorID, currentDate)
    var existingID int
    err = r.db.QueryRow(checkQuery, reflectionID, authorID, currentDate).Scan(&existingID)
    
    if err == sql.ErrNoRows {
        // No tracking record exists for today, insert one
        log.Printf("TrackReflection: No existing tracking found, inserting new record")
        insertQuery := `
            INSERT INTO reflection_tracking (reflection_id, user_id, reflected_date, created_at)
            VALUES ($1, $2, $3, NOW())`
        
        _, err = r.db.Exec(insertQuery, reflectionID, authorID, currentDate)
        if err != nil {
            log.Printf("TrackReflection: Error inserting tracking record: %v", err)
        } else {
            log.Printf("TrackReflection: Successfully inserted tracking record")
        }
        return err
    } else if err != nil {
        log.Printf("TrackReflection: Error checking for existing tracking: %v", err)
        return err
    } else {
        // Tracking record exists for today, delete it (toggle off)
        log.Printf("TrackReflection: Existing tracking found, deleting record")
        deleteQuery := `
            DELETE FROM reflection_tracking 
            WHERE reflection_id = $1 AND user_id = $2 AND reflected_date = $3`
        
        _, err = r.db.Exec(deleteQuery, reflectionID, authorID, currentDate)
        if err != nil {
            log.Printf("TrackReflection: Error deleting tracking record: %v", err)
        } else {
            log.Printf("TrackReflection: Successfully deleted tracking record")
        }
        return err
    }
}

func (r *ReflectionRepository) GetReflectionTrackingAnalytics(userID int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    log.Printf("Getting analytics for user %d from %s to %s", userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    
    query := `
        SELECT 
            rt.reflected_date,
            COUNT(*) as reflection_count,
            COUNT(DISTINCT rt.reflection_id) as unique_reflections
        FROM reflection_tracking rt
        WHERE rt.user_id = $1 
            AND rt.reflected_date >= $2 
            AND rt.reflected_date <= $3
        GROUP BY rt.reflected_date
        ORDER BY rt.reflected_date ASC`
    
    rows, err := r.db.Query(query, userID, startDate, endDate)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var analytics []map[string]interface{}
    for rows.Next() {
        var reflectedDate time.Time
        var reflectionCount, uniqueReflections int
        
        err := rows.Scan(&reflectedDate, &reflectionCount, &uniqueReflections)
        if err != nil {
            return nil, err
        }
        
        log.Printf("Analytics data: date=%s, count=%d, unique=%d", reflectedDate.Format("2006-01-02"), reflectionCount, uniqueReflections)
        analytics = append(analytics, map[string]interface{}{
            "date": reflectedDate.Format("2006-01-02"),
            "reflection_count": reflectionCount,
            "unique_reflections": uniqueReflections,
        })
    }
    
    return analytics, nil
}

func (r *ReflectionRepository) GetReflectionTrackingByCategory(userID int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    query := `
        SELECT 
            DATE(r.created_at) as date,
            c.name as category,
            COUNT(*) as reflection_count
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        WHERE r.author_id = $1 
        AND r.created_at >= $2 
        AND r.created_at <= $3
        GROUP BY DATE(r.created_at), c.name
        ORDER BY date, category`
    
    rows, err := r.db.Query(query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []map[string]interface{}
    for rows.Next() {
        var date time.Time
        var category string
        var reflectionCount int
        err := rows.Scan(&date, &category, &reflectionCount)
        if err != nil {
            return nil, err
        }
        
        results = append(results, map[string]interface{}{
            "date":             date.Format("2006-01-02"),
            "category":         category,
            "reflection_count": reflectionCount,
        })
    }
    
    return results, nil
}

// GetFriendsFeed gets public reflections from friends for the social feed
func (r *ReflectionRepository) GetFriendsFeed(userID int, limit, offset int) ([]models.Reflection, error) {
    fmt.Printf("=== REPOSITORY FEED DEBUG ===\n")
    fmt.Printf("User ID: %d, Limit: %d, Offset: %d\n", userID, limit, offset)
    
    query := `
        SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
               r.reflection_text, r.reflection_detail, r.tags, r.is_private, r.visibility,
               r.created_at, r.updated_at,
               c.name as category_name, c.description as category_desc,
               sc.name as subcategory_name, sc.description as subcategory_desc,
               u.username as author_username, u.first_name, u.last_name
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
        JOIN users u ON r.author_id = u.id
        WHERE r.visibility = 'public'
        AND r.author_id IN (
            SELECT CASE 
                WHEN user_id = $1 THEN friend_id 
                ELSE user_id 
            END
            FROM friendships 
            WHERE (user_id = $1 OR friend_id = $1) 
            AND status = 'accepted'
        )
        AND r.created_at >= NOW() - INTERVAL '14 days'
        ORDER BY r.created_at DESC
        LIMIT $2 OFFSET $3`
    
    fmt.Printf("Executing query with params: userID=%d, limit=%d, offset=%d\n", userID, limit, offset)
    
    rows, err := r.db.Query(query, userID, limit, offset)
    if err != nil {
        fmt.Printf("Database query error: %v\n", err)
        return nil, err
    }
    defer rows.Close()
    
    var reflections []models.Reflection
    rowCount := 0
    for rows.Next() {
        rowCount++
        var reflection models.Reflection
        var categoryName, categoryDesc string
        var subcategoryName, subcategoryDesc sql.NullString
        var firstName, lastName string
        
        err := rows.Scan(
            &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
            &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate, &reflection.Visibility,
            &reflection.CreatedAt, &reflection.UpdatedAt,
            &categoryName, &categoryDesc,
            &subcategoryName, &subcategoryDesc,
            &reflection.AuthorUsername, &firstName, &lastName,
        )
        if err != nil {
            fmt.Printf("Row scan error: %v\n", err)
            return nil, err
        }
        
        reflection.AuthorName = firstName + " " + lastName
        
        // Populate category
        reflection.Category = &models.Category{
            ID:          reflection.CategoryID,
            Name:        categoryName,
            Description: categoryDesc,
        }
        
        // Populate subcategory if exists
        if subcategoryName.Valid {
            reflection.SubCategory = &models.SubCategory{
                ID:          *reflection.SubCategoryID,
                CategoryID:  reflection.CategoryID,
                Name:        subcategoryName.String,
                Description: subcategoryDesc.String,
            }
        }
        
        reflections = append(reflections, reflection)
        fmt.Printf("Processed row %d: reflection ID %d, author %s, visibility %s\n", rowCount, reflection.ID, reflection.AuthorName, reflection.Visibility)
    }
    
    fmt.Printf("Total rows processed: %d\n", rowCount)
    return reflections, nil
}

// GetFriendReflections gets all public reflections from a specific friend
func (r *ReflectionRepository) GetFriendReflections(userID, friendID int, limit, offset int) ([]models.Reflection, error) {
    // First check if they are friends
    var friendshipExists bool
    err := r.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM friendships 
            WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
            AND status = 'accepted'
        )
    `, userID, friendID).Scan(&friendshipExists)
    
    if err != nil {
        return nil, err
    }
    
    if !friendshipExists {
        return nil, fmt.Errorf("users are not friends")
    }
    
    query := `
        SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
               r.reflection_text, r.reflection_detail, r.tags, r.is_private, r.visibility,
               r.created_at, r.updated_at,
               c.name as category_name, c.description as category_desc,
               sc.name as subcategory_name, sc.description as subcategory_desc,
               u.username as author_username, u.first_name, u.last_name
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
        JOIN users u ON r.author_id = u.id
        WHERE r.author_id = $1 AND r.visibility = 'public'
        ORDER BY r.created_at DESC
        LIMIT $2 OFFSET $3`
    
    rows, err := r.db.Query(query, friendID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var reflections []models.Reflection
    for rows.Next() {
        var reflection models.Reflection
        var categoryName, categoryDesc string
        var subcategoryName, subcategoryDesc sql.NullString
        var firstName, lastName string
        
        err := rows.Scan(
            &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
            &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate, &reflection.Visibility,
            &reflection.CreatedAt, &reflection.UpdatedAt,
            &categoryName, &categoryDesc,
            &subcategoryName, &subcategoryDesc,
            &reflection.AuthorUsername, &firstName, &lastName,
        )
        if err != nil {
            return nil, err
        }
        
        reflection.AuthorName = firstName + " " + lastName
        
        // Populate category
        reflection.Category = &models.Category{
            ID:          reflection.CategoryID,
            Name:        categoryName,
            Description: categoryDesc,
        }
        
        // Populate subcategory if exists
        if subcategoryName.Valid {
            reflection.SubCategory = &models.SubCategory{
                ID:          *reflection.SubCategoryID,
                CategoryID:  reflection.CategoryID,
                Name:        subcategoryName.String,
                Description: subcategoryDesc.String,
            }
        }
        
        reflections = append(reflections, reflection)
    }
    
    return reflections, nil
}

// GetFriendReflectionsByUsername gets public reflections from a friend by username
func (r *ReflectionRepository) GetFriendReflectionsByUsername(currentUserID int, friendUsername string, limit, offset int) ([]models.Reflection, error) {
    fmt.Printf("=== REPOSITORY FRIEND BY USERNAME DEBUG ===\n")
    fmt.Printf("Current User ID: %d, Friend Username: %s, Limit: %d, Offset: %d\n", currentUserID, friendUsername, limit, offset)
    
    // First, get the friend's user ID by username
    var friendID int
    err := r.db.QueryRow("SELECT id FROM users WHERE username = $1", friendUsername).Scan(&friendID)
    if err != nil {
        fmt.Printf("Error finding friend by username '%s': %v\n", friendUsername, err)
        if err.Error() == "sql: no rows in result set" {
            return []models.Reflection{}, fmt.Errorf("friend not found: username '%s' does not exist", friendUsername)
        }
        return []models.Reflection{}, fmt.Errorf("database error finding friend: %v", err)
    }
    
    fmt.Printf("Found friend ID: %d for username: %s\n", friendID, friendUsername)
    
    // Check if they are friends
    var friendshipExists bool
    err = r.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM friendships 
            WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
            AND status = 'accepted'
        )
    `, currentUserID, friendID).Scan(&friendshipExists)
    
    if err != nil {
        fmt.Printf("Error checking friendship: %v\n", err)
        return []models.Reflection{}, fmt.Errorf("database error checking friendship: %v", err)
    }
    
    if !friendshipExists {
        fmt.Printf("Users are not friends: currentUserID=%d, friendID=%d\n", currentUserID, friendID)
        return []models.Reflection{}, fmt.Errorf("users are not friends")
    }
    
    fmt.Printf("Friendship confirmed. Getting reflections from friend ID: %d\n", friendID)
    
    // Get public reflections from the friend from the last 14 days
    query := `
        SELECT r.id, r.author_id, r.category_id, r.sub_category_id, r.date, 
               r.reflection_text, r.reflection_detail, r.tags, r.is_private, r.visibility,
               r.created_at, r.updated_at,
               c.name as category_name, c.description as category_desc,
               sc.name as subcategory_name, sc.description as subcategory_desc,
               u.username as author_username, u.first_name, u.last_name
        FROM reflections r
        JOIN categories c ON r.category_id = c.id
        LEFT JOIN sub_categories sc ON r.sub_category_id = sc.id
        JOIN users u ON r.author_id = u.id
        WHERE r.author_id = $1 
        AND r.visibility = 'public'
        AND r.created_at >= NOW() - INTERVAL '14 days'
        ORDER BY r.created_at DESC
        LIMIT $2 OFFSET $3`
    
    fmt.Printf("Executing query with params: friendID=%d, limit=%d, offset=%d\n", friendID, limit, offset)
    
    rows, err := r.db.Query(query, friendID, limit, offset)
    if err != nil {
        fmt.Printf("Database query error: %v\n", err)
        return []models.Reflection{}, fmt.Errorf("database error querying reflections: %v", err)
    }
    defer rows.Close()
    
    var reflections []models.Reflection
    rowCount := 0
    for rows.Next() {
        rowCount++
        var reflection models.Reflection
        var categoryName, categoryDesc string
        var subcategoryName, subcategoryDesc sql.NullString
        var firstName, lastName string
        
        err := rows.Scan(
            &reflection.ID, &reflection.AuthorID, &reflection.CategoryID, &reflection.SubCategoryID, &reflection.Date,
            &reflection.ReflectionText, &reflection.ReflectionDetail, &reflection.Tags, &reflection.IsPrivate, &reflection.Visibility,
            &reflection.CreatedAt, &reflection.UpdatedAt,
            &categoryName, &categoryDesc,
            &subcategoryName, &subcategoryDesc,
            &reflection.AuthorUsername, &firstName, &lastName,
        )
        if err != nil {
            fmt.Printf("Row scan error: %v\n", err)
            return []models.Reflection{}, fmt.Errorf("error scanning reflection row: %v", err)
        }
        
        reflection.AuthorName = firstName + " " + lastName
        
        // Populate category
        reflection.Category = &models.Category{
            ID:          reflection.CategoryID,
            Name:        categoryName,
            Description: categoryDesc,
        }
        
        // Populate subcategory if exists
        if subcategoryName.Valid {
            reflection.SubCategory = &models.SubCategory{
                ID:          *reflection.SubCategoryID,
                CategoryID:  reflection.CategoryID,
                Name:        subcategoryName.String,
                Description: subcategoryDesc.String,
            }
        }
        
        reflections = append(reflections, reflection)
        fmt.Printf("Processed row %d: reflection ID %d, author %s, visibility %s\n", rowCount, reflection.ID, reflection.AuthorName, reflection.Visibility)
    }
    
    fmt.Printf("Total rows processed: %d\n", rowCount)
    fmt.Printf("Returning %d reflections\n", len(reflections))
    return reflections, nil
}