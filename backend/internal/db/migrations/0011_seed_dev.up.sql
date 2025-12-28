-- seed dev user
insert into app_user (email, display_name, password_hash, is_active)
values ('dev@local.seed', 'Dev User', 'v1:bcrypt:$2a$12$n180ANBjuXfZrr.hWFZXjukiDZuQ1Kw6yauaIrEHriMjempCALOB2', true);

-- seed libraries
insert into library (name, type, root_path, enabled, "default")
values ('Main Movie Library', 'movie', '/data/movies', true, true);

insert into library (name, type, root_path, enabled, "default")
values ('Main Series Library', 'series', '/data/tv', true, true);

-- seed name templates
insert into name_template (name, type, template, "default")
values ('Main Movie Template', 'movie', '{{.Title}} ({{.Year}}) [{{clean .Quality.Resolution}} {{clean .Quality.Codec}}]', true);

insert into name_template (name, type, template, "default")
values ('Main Series Template', 'series', '{{.Title}} - S{{.Season}}E{{.Episode}} - {{.EpisodeTitle}} ({{.Year}}) [{{clean .Quality.Resolution}} {{clean .Quality.Codec}}]', true);

-- seed downloaders
insert into downloader (name, type, protocol, url, username, password, enabled, "default")
values ('Main Downloader', 'qbittorrent', 'torrent', 'http://172.16.10.22:8485', 'admin', 'admin', true, true);

