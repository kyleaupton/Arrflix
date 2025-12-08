-- add default name templates
insert into name_template (name, type, template, "default")
values ('Default Movie Template', 'movie', '{Title} ({Year}) {Quality} {Resolution} {Extension}', true);
insert into name_template (name, type, template, "default")
values ('Default Series Template', 'series', '{Title} - S{Season:00}E{Episode:00} - {EpisodeTitle} ({Year}) {Quality} {Resolution} {Extension}', true);

-- add default movie policy
insert into policy (name, description, enabled, priority)
values ('Default Movie Policy', 'Default policy for movies', true, 0);

-- add default movie rule
insert into rule (policy_id, left_operand, operator, right_operand)
values (
    (select id from policy where name = 'Default Movie Policy'),
    'torrent.type',
    '==',
    'movie'
);

-- set downloader for default movie policy
insert into action (policy_id, type, value, "order")
values (
    (select id from policy where name = 'Default Movie Policy'), 
    'set_downloader',
    'qbittorrent',
    1
);

-- set library for default movie policy
insert into action (policy_id, type, value, "order")
values (
    (select id from policy where name = 'Default Movie Policy'), 
    'set_library',
    (select id from library where name = 'Test Movie Library'),
    2
);

-- set name template for default movie policy
insert into action (policy_id, type, value, "order")
values (
    (select id from policy where name = 'Default Movie Policy'), 
    'set_name_template',
    (select id from name_template where name = 'Default Movie Template'),
    3
);