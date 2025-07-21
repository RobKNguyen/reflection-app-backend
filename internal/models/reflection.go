// Updated internal/models/reflection.go
package models

import (
    "time"
    "database/sql/driver"
    "encoding/json"
    "errors"
    "strings"
    "log"
)

// Custom type for string arrays in PostgreSQL
type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
    if len(sa) == 0 {
        return "{}", nil
    }
    
    // Convert to PostgreSQL array format: {"item1","item2","item3"}
    result := "{"
    for i, item := range sa {
        if i > 0 {
            result += ","
        }
        // Escape quotes and backslashes
        escaped := strings.ReplaceAll(item, `\`, `\\`)
        escaped = strings.ReplaceAll(escaped, `"`, `\"`)
        result += `"` + escaped + `"`
    }
    result += "}"
    
    return result, nil
}

func (sa *StringArray) Scan(value interface{}) error {
    if value == nil {
        *sa = StringArray{}
        return nil
    }
    
    switch v := value.(type) {
    case []byte:
        // Handle PostgreSQL array format: {"item1","item2","item3"}
        str := string(v)
        if str == "{}" {
            *sa = StringArray{}
            return nil
        }
        
        // Remove the outer braces
        str = strings.Trim(str, "{}")
        if str == "" {
            *sa = StringArray{}
            return nil
        }
        
        // Split by comma and unquote
        parts := strings.Split(str, ",")
        result := make(StringArray, 0, len(parts))
        for _, part := range parts {
            part = strings.TrimSpace(part)
            if part != "" {
                // Remove quotes
                part = strings.Trim(part, `"`)
                result = append(result, part)
            }
        }
        *sa = result
        return nil
        
    case string:
        // Handle string format
        if v == "{}" {
            *sa = StringArray{}
            return nil
        }
        return sa.Scan([]byte(v))
        
    default:
        return errors.New("cannot scan into StringArray")
    }
}

// Custom date type for JSON serialization
type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
    var dateStr string
    if err := json.Unmarshal(data, &dateStr); err != nil {
        return err
    }
    
    parsedDate, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
    if err != nil {
        return err
    }
    
    *d = DateOnly(parsedDate)
    return nil
}

func (d DateOnly) Time() time.Time {
    return time.Time(d)
}

type Reflection struct {
    ID               int         `json:"id" db:"id"`
    AuthorID         int         `json:"author_id" db:"author_id"`
    CategoryID       int         `json:"category_id" db:"category_id"`
    SubCategoryID    *int        `json:"sub_category_id" db:"sub_category_id"` // nullable
    Date             DateOnly    `json:"date" db:"date"`
    ReflectionText   string      `json:"reflection_text" db:"reflection_text"`
    ReflectionDetail string      `json:"reflection_detail" db:"reflection_detail"`
    Tags             StringArray `json:"tags" db:"tags"`
    IsPrivate        bool        `json:"is_private" db:"is_private"`
    Visibility       string      `json:"visibility" db:"visibility"` // private, public
    CreatedAt        time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
    
    // Related data (populated when needed)
    Category    *Category    `json:"category,omitempty"`
    SubCategory *SubCategory `json:"sub_category,omitempty"`
    Actions     []ActionItem `json:"actions,omitempty"`
    
    // Tracking data
    ReflectionCount int  `json:"reflection_count"`
    ReflectedToday  bool `json:"reflected_today"`
    
    // Social data (populated when needed)
    AuthorUsername string `json:"author_username,omitempty"`
    AuthorName     string `json:"author_name,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Reflection
func (r *Reflection) UnmarshalJSON(data []byte) error {
    type Alias Reflection
    aux := &struct {
        Date string `json:"date"`
        *Alias
    }{
        Alias: (*Alias)(r),
    }
    
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    
    // Parse the date string if it's not empty
    if aux.Date != "" {
        log.Printf("UnmarshalJSON: Received date string: '%s'", aux.Date)
        // Parse the date string as a local date (not UTC)
        parsedDate, err := time.ParseInLocation("2006-01-02", aux.Date, time.Local)
        if err != nil {
            log.Printf("UnmarshalJSON: Error parsing date '%s': %v", aux.Date, err)
            return err
        }
        log.Printf("UnmarshalJSON: Parsed date: %s (Local: %s, UTC: %s)", 
            parsedDate.Format("2006-01-02"), 
            parsedDate.Local().Format("2006-01-02 15:04:05 MST"),
            parsedDate.UTC().Format("2006-01-02 15:04:05 UTC"))
        r.Date = DateOnly(parsedDate)
    } else {
        // If no date provided, use current date
        currentDate := time.Now().Truncate(24 * time.Hour)
        log.Printf("UnmarshalJSON: No date provided, using current date: %s", currentDate.Format("2006-01-02"))
        r.Date = DateOnly(currentDate)
    }
    
    return nil
}