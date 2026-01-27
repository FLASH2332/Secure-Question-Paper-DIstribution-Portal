package acl

import (
    "database/sql"
    "fmt"
)

// Permission represents what a role can do with an object
type Permission struct {
    Role       string
    ObjectType string
    CanCreate  bool
    CanRead    bool
    CanUpdate  bool
    CanDelete  bool
    CanEncrypt bool
    CanDecrypt bool
}

// GetPermissions retrieves permissions for a role and object type
func GetPermissions(db *sql.DB, role, objectType string) (*Permission, error) {
    var perm Permission
    
    query := `
        SELECT role, object_type, can_create, can_read, can_update, can_delete, can_encrypt, can_decrypt
        FROM access_control
        WHERE role = ? AND object_type = ?
    `
    
    err := db.QueryRow(query, role, objectType).Scan(
        &perm.Role,
        &perm.ObjectType,
        &perm.CanCreate,
        &perm.CanRead,
        &perm.CanUpdate,
        &perm.CanDelete,
        &perm.CanEncrypt,
        &perm.CanDecrypt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("no permissions found for role %s on %s", role, objectType)
    } else if err != nil {
        return nil, fmt.Errorf("failed to get permissions: %w", err)
    }
    
    return &perm, nil
}

// CheckPermission verifies if a role has a specific permission
func CheckPermission(db *sql.DB, role, objectType, action string) (bool, error) {
    perm, err := GetPermissions(db, role, objectType)
    if err != nil {
        return false, err
    }
    
    switch action {
    case "create":
        return perm.CanCreate, nil
    case "read":
        return perm.CanRead, nil
    case "update":
        return perm.CanUpdate, nil
    case "delete":
        return perm.CanDelete, nil
    case "encrypt":
        return perm.CanEncrypt, nil
    case "decrypt":
        return perm.CanDecrypt, nil
    default:
        return false, fmt.Errorf("unknown action: %s", action)
    }
}

// GetAllPermissions retrieves all permissions for a role
func GetAllPermissions(db *sql.DB, role string) ([]Permission, error) {
    query := `
        SELECT role, object_type, can_create, can_read, can_update, can_delete, can_encrypt, can_decrypt
        FROM access_control
        WHERE role = ?
    `
    
    rows, err := db.Query(query, role)
    if err != nil {
        return nil, fmt.Errorf("failed to get permissions: %w", err)
    }
    defer rows.Close()
    
    var permissions []Permission
    for rows.Next() {
        var perm Permission
        err := rows.Scan(
            &perm.Role,
            &perm.ObjectType,
            &perm.CanCreate,
            &perm.CanRead,
            &perm.CanUpdate,
            &perm.CanDelete,
            &perm.CanEncrypt,
            &perm.CanDecrypt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan permission: %w", err)
        }
        permissions = append(permissions, perm)
    }
    
    return permissions, nil
}