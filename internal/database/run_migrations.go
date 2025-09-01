package database

import (
    "collabhub-music-backend/internal/models"
    "fmt"

    "gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
    err := db.AutoMigrate(
        &models.User{},
        &models.Project{},
        &models.ProjectCollaborator{},
        &models.Branch{},
        &models.File{},
        &models.FileVersion{},
        &models.AudioMetadata{},
    )
    if err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    return nil
}