-- Media queries

-- name: ListMediaItems :many
select * from media_item
order by created_at desc;

-- name: GetMediaItem :one
select * from media_item
where id = $1;

-- name: GetMediaItemByTmdbID :one
select * from media_item
where tmdb_id = $1;

-- name: CreateMediaItem :one
insert into media_item (library_id, type, title, year, tmdb_id)
values (sqlc.arg(library_id), sqlc.arg(type), sqlc.arg(title), sqlc.arg(year), sqlc.arg(tmdb_id))
returning *;

-- name: UpdateMediaItem :one
update media_item
set title = $2,
    year = $3,
    tmdb_id = $4,
    updated_at = now()
where id = $1
returning *;

-- name: DeleteMediaItem :exec
delete from media_item where id = $1;

-- Seasons

-- name: ListSeasonsForMedia :many
select * from media_season
where media_item_id = $1
order by season_number asc;

-- name: UpsertSeason :one
insert into media_season (media_item_id, season_number, air_date)
values (sqlc.arg(media_item_id), sqlc.arg(season_number), sqlc.arg(air_date))
on conflict (media_item_id, season_number)
do update set air_date = excluded.air_date
returning *;

-- Episodes

-- name: ListEpisodesForSeason :many
select * from media_episode
where season_id = $1
order by episode_number asc;

-- name: UpsertEpisode :one
insert into media_episode (season_id, episode_number, title, air_date, tmdb_id, tvdb_id)
values (sqlc.arg(season_id), sqlc.arg(episode_number), sqlc.arg(title), sqlc.arg(air_date), sqlc.arg(tmdb_id), sqlc.arg(tvdb_id))
on conflict (season_id, episode_number)
do update set title = excluded.title,
              air_date = excluded.air_date,
              tmdb_id = excluded.tmdb_id,
              tvdb_id = excluded.tvdb_id
returning *;

-- Files

-- name: GetMediaFileByPath :one
select * from media_file where path = $1;

-- name: CreateMediaFile :one
insert into media_file (media_item_id, season_id, episode_id, path)
values (sqlc.arg(media_item_id), sqlc.arg(season_id), sqlc.arg(episode_id), sqlc.arg(path))
returning *;

-- name: DeleteMediaFile :exec
delete from media_file where id = $1;




