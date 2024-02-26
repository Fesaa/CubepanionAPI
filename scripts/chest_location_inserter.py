import asyncio
import re

import asyncpg
import toml

class ChestLocation:

    def __init__(self, x: int, y: int, z: int) -> None:
        self.x = x
        self.y = y
        self.z = z

    def __str__(self) -> str:
        return f"ChestLocation({self.x}, {self.y}, {self.z})"


async def main():
    config = toml.load("config.toml")
    db_url = config["database_url"]

    pool = await asyncpg.create_pool(db_url)

    path = input("Give a path to the .txt file containing chest data: ")
    if not path.endswith(".txt"):
        print("Adding .txt to path")
        path += ".txt"

    season = input("Give the season you're submitting for: ")
    running = True if input("Running season? ") in ("true", "t", "1") else False

    data = open(path, encoding="utf-8")
    chest_locations: list[ChestLocation] = []
    for line in data.readlines():
        match = re.match(".* X: (-?)(\d{1,9}).*, Y: (-?)(\d{1,9}).*, Z: (-?)(\d{1,9}).*", line)
        if match != None and len(match.groups()) == 6:

            x = -int(match.groups()[1]) if match.groups()[0] == "-" else int(match.groups()[1])
            y = -int(match.groups()[3]) if match.groups()[2] == "-" else int(match.groups()[3])
            z = -int(match.groups()[5]) if match.groups()[4] == "-" else int(match.groups()[5])
            chest_locations.append(ChestLocation(x, y, z))

    async with pool.acquire() as con:
            con: asyncpg.Connection
            await con.execute("INSERT INTO seasons (season_name, running) VALUES ($1, $2) ON CONFLICT (season_name) DO NOTHING;", season, running)
            values = ",".join(f"('{season}', {loc.x}, {loc.y}, {loc.z})" for loc in chest_locations)
            query = f"INSERT INTO chest_locations (season_name, x, y, z) VALUES {values};"
            #print(query)
            await con.execute(query)

if __name__ == "__main__":
    asyncio.run(main())