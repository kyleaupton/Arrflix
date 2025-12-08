-- seed dev user
insert into app_user (email, display_name, password_hash, is_active)
values ('dev@local.seed', 'Dev User', 'v1:bcrypt:$2a$12$n180ANBjuXfZrr.hWFZXjukiDZuQ1Kw6yauaIrEHriMjempCALOB2', true);

-- seed libraries
insert into library (name, type, root_path, enabled, "default")
values ('Main Movie Library', 'movie', '/libraries/movies', true, true);

insert into library (name, type, root_path, enabled, "default")
values ('Main Series Library', 'series', '/libraries/tv', true, false);

-- seed name templates
insert into name_template (name, type, template, "default")
values ('Main Movie Template', 'movie', '{Title} ({Year}) {Quality} {Resolution} {Extension}', true);

insert into name_template (name, type, template, "default")
values ('Main Series Template', 'series', '{Title} - S{Season:00}E{Episode:00} - {EpisodeTitle} ({Year}) {Quality} {Resolution} {Extension}', true);

-- seed downloaders
insert into downloader (name, type, protocol, url, username, password, enabled, "default")
values ('Main Downloader', 'qbittorrent', 'torrent', 'http://localhost:8080', 'admin', 'admin', true, true);

