-- seed dev user
insert into app_user (email, display_name, password_hash, is_active)
values ('dev@local.seed', 'Dev User', 'v1:bcrypt:$2a$12$n180ANBjuXfZrr.hWFZXjukiDZuQ1Kw6yauaIrEHriMjempCALOB2', true);

-- seed libraries
insert into library (name, type, root_path, enabled, "default")
values ('Main Movie Library', 'movie', '/data/movies', true, true);

insert into library (name, type, root_path, enabled, "default")
values ('Main Series Library', 'series', '/data/tv', true, true);

-- seed name templates
insert into name_template (name, type, template, series_show_template, series_season_template, "default")
values ('Main Movie Template', 'movie', '{{.Media.CleanTitle}} ({{.Media.Year}}) {tmdb-{{.Media.TmdbID}}} {{if .Release.Edition}}{edition-{{.Release.Edition}}}{{end}} [{{.Quality.Full}}][{{.MediaInfo.AudioCodec}} {{.MediaInfo.AudioChannels}}][{{.MediaInfo.VideoCodec}}]{{if .Release.ReleaseGroup}}-{{.Release.ReleaseGroup}}{{end}}', null, null, true);

insert into name_template (name, type, template, series_show_template, series_season_template, "default")
values ('Main Series Template', 'series', '{{.Media.Title}} - S{{.Media.Season}}E{{.Media.Episode}} - {{clean .Media.EpisodeTitle}} ({{.Media.Year}}) [{{.Quality.Full}}][{{.MediaInfo.AudioCodec}} {{.MediaInfo.AudioChannels}}][{{.MediaInfo.VideoCodec}}]{{if .Release.ReleaseGroup}}-{{.Release.ReleaseGroup}}{{end}}', '{{.Media.Title}} ({{.Media.Year}})', 'Season {{.Media.Season}}', true);

-- seed downloaders
insert into downloader (name, type, protocol, url, username, password, enabled, "default")
values ('Main Downloader', 'qbittorrent', 'torrent', 'http://172.16.10.22:8485', 'admin', 'admin', true, true);

