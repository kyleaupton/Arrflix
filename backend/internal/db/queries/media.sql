-- Media queries

-- name: ListMediaItems :many
select * from media_item
order by created_at desc;

-- name: GetMediaItem :one
select * from media_item
where id = $1;

-- name: GetMediaItemByTmdbIDAndType :one
select * from media_item
where tmdb_id = $1 and type = $2;

-- name: GetMediaItemByTmdbID :one
select * from media_item
where tmdb_id = $1;

-- name: CreateMediaItem :one
insert into media_item (type, title, year, tmdb_id)
values (sqlc.arg(type), sqlc.arg(title), sqlc.arg(year), sqlc.arg(tmdb_id))
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

-- name: GetSeason :one
select * from media_season
where id = $1;

-- name: GetSeasonByNumber :one
select * from media_season
where media_item_id = $1 and season_number = $2;

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

-- name: GetEpisode :one
select * from media_episode
where id = $1;

-- name: GetEpisodeByNumber :one
select * from media_episode
where season_id = $1 and episode_number = $2;

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

-- name: GetMediaFileByLibraryAndPath :one
select * from media_file where library_id = $1 and path = $2;

-- name: CreateMediaFile :one
insert into media_file (library_id, media_item_id, season_id, episode_id, path, status)
values (sqlc.arg(library_id), sqlc.arg(media_item_id), sqlc.arg(season_id), sqlc.arg(episode_id), sqlc.arg(path), coalesce(sqlc.arg(status), 'available'))
returning *;

-- name: DeleteMediaFile :exec
delete from media_file where id = $1;

-- name: ListMediaFilesForItem :many
select
  mf.id,
  mf.library_id,
  mf.media_item_id,
  mf.season_id,
  mf.episode_id,
  mf.path,
  mf.status,
  mf.added_at,
  ms.season_number,
  me.episode_number
from media_file mf
left join media_season ms on mf.season_id = ms.id
left join media_episode me on mf.episode_id = me.id
where mf.media_item_id = $1
order by mf.added_at desc;

-- name: ListEpisodeAvailabilityForSeries :many
select
  ms.season_number,
  me.episode_number,
  me.id as episode_id,
  me.title,
  me.air_date,
  mf.id as file_id,
  mf.library_id,
  mf.status
from media_episode me
join media_season ms on me.season_id = ms.id
join media_item mi on ms.media_item_id = mi.id
left join media_file mf on mf.episode_id = me.id
where mi.id = $1
order by ms.season_number, me.episode_number;




