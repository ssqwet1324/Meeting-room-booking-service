-- Уникальность пары (комната, начало слота) нужна для ON CONFLICT
CREATE UNIQUE INDEX IF NOT EXISTS slots_room_id_start_time_unique ON slots (room_id, start_time);
