import asyncio
import asyncpg
import json
import toml

class Generator:

    def __init__(self, unique_name: str, ordering: int, gen_type: str, gen_location: str, level: int, count: int) -> None:
        self.unique_name = unique_name
        self.ordering = ordering
        self.gen_type = gen_type
        self.gen_location = gen_location
        self.level = level
        self.count = count


    def __str__(self) -> str:
        return f"Generator({self.unique_name}, {self.ordering}, {self.gen_type}, {self.gen_location}, {self.level}, {self.count})"

class EggWarsMap:

    def __init__(self, unique_name: str, map_name: str, team_size: int, build_limit: int, colours: str, generators: list[dict]) -> None:
        self.unique_name = unique_name
        self.map_name = map_name
        self.team_size = team_size
        self.build_limit = build_limit
        self.colours = colours
        self.generators = [Generator(gen["unique_name"], gen["ordering"], gen["gen_type"], gen["gen_location"], gen["level"], gen["count"]) for gen in generators]

    def __str__(self) -> str:
        return f"EggWarsMap({self.unique_name}, {self.map_name}, {self.team_size}, {self.build_limit}, {self.colours}, {self.generators})"



async def main():
    config = toml.load("config.toml")
    db_url = config["database_url"]

    pool = await asyncpg.create_pool(db_url)

    with open("eggwars_maps.json") as f:
        maps_data = json.load(f)
        maps = [EggWarsMap(m["unique_name"], m["map_name"], m["team_size"], m["build_limit"], m["colours"], m["generators"]) for _, m in maps_data.items()]

    eggwars_maps_query = "INSERT INTO eggwars_maps (unique_name, map_name, team_size, build_limit, colours) VALUES "
    gen_layout_query = "INSERT INTO gen_layout (unique_name, ordering, gen_type, gen_location, level, count) VALUES"

    for eggwars_map in maps:
        layout = eggwars_map.generators
        eggwars_maps_query += f"('{eggwars_map.unique_name}', '{eggwars_map.map_name}', {eggwars_map.team_size}, {eggwars_map.build_limit}, '{eggwars_map.colours}'), "
        for gen in layout:
            gen_layout_query += f"('{gen.unique_name}', {gen.ordering}, '{gen.gen_type}', '{gen.gen_location}', {gen.level}, {gen.count}), "

    eggwars_maps_query = eggwars_maps_query.removesuffix(", ") + " ON CONFLICT (unique_name) DO NOTHING;"
    gen_layout_query = gen_layout_query.removesuffix(", ")+ " ON CONFLICT (unique_name, ordering) DO NOTHING;"

    async with pool.acquire() as con:
        await con.execute(eggwars_maps_query)
        await con.execute(gen_layout_query)

    

if __name__ == "__main__":
    asyncio.run(main())

