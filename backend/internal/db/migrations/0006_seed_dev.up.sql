-- seed dev user
insert into app_user (email, display_name, password_hash, is_active)
values ('dev@local.seed', 'Dev User', 'v1:bcrypt:$2a$12$n180ANBjuXfZrr.hWFZXjukiDZuQ1Kw6yauaIrEHriMjempCALOB2', true);

-- seed dev libraries
insert into library (name, type, root_path, enabled)
values ('Test Movie Library', 'movie', '/mnt/media-pipeline/testing/movies', true);

insert into library (name, type, root_path, enabled)
values ('Test Series Library', 'series', '/mnt/media-pipeline/testing/tv', true);
