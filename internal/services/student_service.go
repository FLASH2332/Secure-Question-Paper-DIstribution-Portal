package services

import (
    "database/sql"
    "fmt"
    
    "github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/acl"
    "github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

// StudentService handles student operations
type StudentService struct {
    DB   *sql.DB
    User *models.User
}

// NewStudentService creates a new student service
func NewStudentService(db *sql.DB, user *models.User) *StudentService {
    return &StudentService{
        DB:   db,
        User: user,
    }
}

// CanViewExamSchedule checks if student can view exam schedule
func (s *StudentService) CanViewExamSchedule() error {
    return acl.EnforcePermission(s.DB, s.User, "ExamSession", "read", nil)
}

// CanAccessPaper checks if student can access question paper (should fail)
func (s *StudentService) CanAccessPaper() error {
    return acl.EnforcePermission(s.DB, s.User, "QuestionPaper", "read", nil)
}

// GetExamSchedule retrieves upcoming exam sessions
func (s *StudentService) GetExamSchedule() ([]models.ExamSession, error) {
    // First check permission
    if err := s.CanViewExamSchedule(); err != nil {
        return nil, err
    }
    
    query := `
        SELECT es.id, es.session_name, es.scheduled_time, es.duration_minutes, 
               es.status, qp.title as paper_title, qp.subject
        FROM exam_sessions es
        JOIN question_papers qp ON es.paper_id = qp.id
        WHERE es.status IN ('scheduled', 'active')
        ORDER BY es.scheduled_time ASC
    `
    
    rows, err := s.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to get exam schedule: %w", err)
    }
    defer rows.Close()
    
    var sessions []models.ExamSession
    for rows.Next() {
        var session models.ExamSession
        var paperTitle, subject string
        
        err := rows.Scan(
            &session.ID,
            &session.SessionName,
            &session.ScheduledTime,
            &session.DurationMinutes,
            &session.Status,
            &paperTitle,
            &subject,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan session: %w", err)
        }
        
        sessions = append(sessions, session)
    }
    
    return sessions, nil
}