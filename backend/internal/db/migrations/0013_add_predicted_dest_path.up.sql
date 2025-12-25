-- Add predicted_dest_path column to download_job table
-- This stores the pre-calculated destination path with .{ext} placeholder
ALTER TABLE download_job ADD COLUMN predicted_dest_path text;

