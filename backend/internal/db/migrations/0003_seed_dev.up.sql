-- seed dev user

insert into app_user (email, display_name, password_hash, is_active)
values ('dev@local.host', 'Dev User', 'v1:bcrypt:$2a$12$n180ANBjuXfZrr.hWFZXjukiDZuQ1Kw6yauaIrEHriMjempCALOB2', true);
