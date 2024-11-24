CREATE TABLE IF NOT EXISTS pokemon (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    dex_id INTEGER NOT NULL,
    up_votes INTEGER NOT NULL,
    down_votes INTEGER NOT NULL,
    inserted_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pokemon_dex_id ON pokemon(dex_id);