-- Add movie_dir_template column for movie directory naming
ALTER TABLE name_template ADD COLUMN movie_dir_template text;

-- Backfill existing movie templates with default value
UPDATE name_template
SET movie_dir_template = '{{.Media.CleanTitle}} ({{.Media.Year}}) {tmdb-{{.Media.TmdbID}}}'
WHERE type = 'movie';
