package repository

import (
    "context"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
)

// UserRepository defines the methods for user data access
type UserRepository interface {
    CreateUser(ctx context.Context, user *models.User) error
    GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    GetUserByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    UpdateUser(ctx context.Context, user *models.User) error
    DeleteUser(ctx context.Context, id uuid.UUID) error
    ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
    GetUsersByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*models.User, error)
}

// ProjectRepository defines the methods for project data access
type ProjectRepository interface {
    CreateProject(ctx context.Context, project *models.Project) error
    GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
    UpdateProject(ctx context.Context, project *models.Project) error
    DeleteProject(ctx context.Context, id uuid.UUID) error
    GetProjectsByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*models.Project, error)
    GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error)
    ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, error)
    SearchProjectsByName(ctx context.Context, name string) ([]*models.Project, error)
}

// OrganizationRepository defines the methods for organization data access
type OrganizationRepository interface {
    CreateOrganization(ctx context.Context, org *models.Organization) error
    GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
    UpdateOrganization(ctx context.Context, org *models.Organization) error
    DeleteOrganization(ctx context.Context, id uuid.UUID) error
    GetOrganizationsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error)
    ListOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error)
    GetOrganizationByName(ctx context.Context, name string) (*models.Organization, error)
    AddUserToOrganization(ctx context.Context, orgID, userID uuid.UUID) error
    RemoveUserFromOrganization(ctx context.Context, orgID, userID uuid.UUID) error
}

// CollaborationRepository defines the methods for collaboration data access
type CollaborationRepository interface {
    CreateCollaboration(ctx context.Context, collaboration *models.Collaboration) error
    GetCollaborationByID(ctx context.Context, id uuid.UUID) (*models.Collaboration, error)
    UpdateCollaboration(ctx context.Context, collaboration *models.Collaboration) error
    DeleteCollaboration(ctx context.Context, id uuid.UUID) error
    GetCollaborationsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Collaboration, error)
    GetCollaborationsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Collaboration, error)
}

// FileRepository defines the methods for file data access
type FileRepository interface {
    CreateFile(ctx context.Context, file *models.File) error
    GetFileByID(ctx context.Context, id uuid.UUID) (*models.File, error)
    UpdateFile(ctx context.Context, file *models.File) error
    DeleteFile(ctx context.Context, id uuid.UUID) error
    GetFilesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.File, error)
    GetFilesByUserID(ctx context.Context, userID uuid.UUID) ([]*models.File, error)
    GetFileVersions(ctx context.Context, fileID uuid.UUID) ([]*models.File, error)
}

// TrackRepository defines the methods for track/song data access
type TrackRepository interface {
    CreateTrack(ctx context.Context, track *models.Track) error
    GetTrackByID(ctx context.Context, id uuid.UUID) (*models.Track, error)
    UpdateTrack(ctx context.Context, track *models.Track) error
    DeleteTrack(ctx context.Context, id uuid.UUID) error
    GetTracksByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Track, error)
    GetTracksByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Track, error)
    SearchTracksByName(ctx context.Context, name string) ([]*models.Track, error)
    GetTracksByGenre(ctx context.Context, genre string) ([]*models.Track, error)
    ListTracks(ctx context.Context, limit, offset int) ([]*models.Track, error)
}

// AlbumRepository defines the methods for album data access
type AlbumRepository interface {
    CreateAlbum(ctx context.Context, album *models.Album) error
    GetAlbumByID(ctx context.Context, id uuid.UUID) (*models.Album, error)
    UpdateAlbum(ctx context.Context, album *models.Album) error
    DeleteAlbum(ctx context.Context, id uuid.UUID) error
    GetAlbumsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Album, error)
    GetAlbumsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Album, error)
    AddTrackToAlbum(ctx context.Context, albumID, trackID uuid.UUID, position int) error
    RemoveTrackFromAlbum(ctx context.Context, albumID, trackID uuid.UUID) error
    GetTracksInAlbum(ctx context.Context, albumID uuid.UUID) ([]*models.Track, error)
    ListAlbums(ctx context.Context, limit, offset int) ([]*models.Album, error)
}

// CommentRepository defines the methods for comment data access
type CommentRepository interface {
    CreateComment(ctx context.Context, comment *models.Comment) error
    GetCommentByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
    UpdateComment(ctx context.Context, comment *models.Comment) error
    DeleteComment(ctx context.Context, id uuid.UUID) error
    GetCommentsByTrackID(ctx context.Context, trackID uuid.UUID) ([]*models.Comment, error)
    GetCommentsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Comment, error)
    GetCommentsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Comment, error)
    GetReplies(ctx context.Context, parentCommentID uuid.UUID) ([]*models.Comment, error)
}

// LyricsRepository defines the methods for lyrics data access
type LyricsRepository interface {
    CreateLyrics(ctx context.Context, lyrics *models.Lyrics) error
    GetLyricsByID(ctx context.Context, id uuid.UUID) (*models.Lyrics, error)
    UpdateLyrics(ctx context.Context, lyrics *models.Lyrics) error
    DeleteLyrics(ctx context.Context, id uuid.UUID) error
    GetLyricsByTrackID(ctx context.Context, trackID uuid.UUID) ([]*models.Lyrics, error)
    GetLyricsVersions(ctx context.Context, lyricsID uuid.UUID) ([]*models.Lyrics, error)
}

// InstrumentRepository defines the methods for instrument data access
type InstrumentRepository interface {
    CreateInstrument(ctx context.Context, instrument *models.Instrument) error
    GetInstrumentByID(ctx context.Context, id uuid.UUID) (*models.Instrument, error)
    UpdateInstrument(ctx context.Context, instrument *models.Instrument) error
    DeleteInstrument(ctx context.Context, id uuid.UUID) error
    ListInstruments(ctx context.Context, limit, offset int) ([]*models.Instrument, error)
    GetInstrumentsByCategory(ctx context.Context, category string) ([]*models.Instrument, error)
    SearchInstrumentsByName(ctx context.Context, name string) ([]*models.Instrument, error)
}

// GenreRepository defines the methods for genre data access
type GenreRepository interface {
    CreateGenre(ctx context.Context, genre *models.Genre) error
    GetGenreByID(ctx context.Context, id uuid.UUID) (*models.Genre, error)
    UpdateGenre(ctx context.Context, genre *models.Genre) error
    DeleteGenre(ctx context.Context, id uuid.UUID) error
    ListGenres(ctx context.Context, limit, offset int) ([]*models.Genre, error)
    GetGenreByName(ctx context.Context, name string) (*models.Genre, error)
}

// BranchRepository defines the methods for musical project branches (like Git branches)
type BranchRepository interface {
    CreateBranch(ctx context.Context, branch *models.Branch) error
    GetBranchByID(ctx context.Context, id uuid.UUID) (*models.Branch, error)
    UpdateBranch(ctx context.Context, branch *models.Branch) error
    DeleteBranch(ctx context.Context, id uuid.UUID) error
    GetBranchesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Branch, error)
    MergeBranch(ctx context.Context, sourceBranchID, targetBranchID uuid.UUID) error
    GetBranchHistory(ctx context.Context, branchID uuid.UUID) ([]*models.Commit, error)
}

// CommitRepository defines the methods for musical commit data access
type CommitRepository interface {
    CreateCommit(ctx context.Context, commit *models.Commit) error
    GetCommitByID(ctx context.Context, id uuid.UUID) (*models.Commit, error)
    GetCommitsByBranchID(ctx context.Context, branchID uuid.UUID) ([]*models.Commit, error)
    GetCommitsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Commit, error)
    GetCommitsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Commit, error)
}

// PlaylistRepository defines the methods for playlist data access
type PlaylistRepository interface {
    CreatePlaylist(ctx context.Context, playlist *models.Playlist) error
    GetPlaylistByID(ctx context.Context, id uuid.UUID) (*models.Playlist, error)
    UpdatePlaylist(ctx context.Context, playlist *models.Playlist) error
    DeletePlaylist(ctx context.Context, id uuid.UUID) error
    GetPlaylistsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Playlist, error)
    AddTrackToPlaylist(ctx context.Context, playlistID, trackID uuid.UUID, position int) error
    RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID uuid.UUID) error
    GetTracksInPlaylist(ctx context.Context, playlistID uuid.UUID) ([]*models.Track, error)
}

// ReviewRepository defines the methods for musical review data access
type ReviewRepository interface {
    CreateReview(ctx context.Context, review *models.Review) error
    GetReviewByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
    UpdateReview(ctx context.Context, review *models.Review) error
    DeleteReview(ctx context.Context, id uuid.UUID) error
    GetReviewsByTrackID(ctx context.Context, trackID uuid.UUID) ([]*models.Review, error)
    GetReviewsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Review, error)
    GetReviewsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Review, error)
}

// TagRepository defines the methods for tag data access
type TagRepository interface {
    CreateTag(ctx context.Context, tag *models.Tag) error
    GetTagByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
    UpdateTag(ctx context.Context, tag *models.Tag) error
    DeleteTag(ctx context.Context, id uuid.UUID) error
    GetTagsByName(ctx context.Context, name string) ([]*models.Tag, error)
    GetTagsByTrackID(ctx context.Context, trackID uuid.UUID) ([]*models.Tag, error)
    AddTagToTrack(ctx context.Context, tagID, trackID uuid.UUID) error
    RemoveTagFromTrack(ctx context.Context, tagID, trackID uuid.UUID) error
    ListTags(ctx context.Context, limit, offset int) ([]*models.Tag, error)
}