import re

class ChestLocation:

    def __init__(self, x: int, y: int, z: int) -> None:
        self.x = x
        self.y = y
        self.z = z

    def __str__(self) -> str:
        return f"ChestLocation({self.x}, {self.y}, {self.z})"


def main():
    path = input("Give a path to the .txt file containing chest data: ")
    season = input("Give the season you're submitting for: ")
    running = True if input("Running season? ") in ("true", "t", "1") else False

    data = open(path, encoding="utf-8")
    chest_locations: list[ChestLocation] = []
    for line in data.readlines():
        line = line.replace('§b', '').replace('§6','')
        match = re.match(".* X: (-?)(\d{1,9}).*, Y: (-?)(\d{1,9}).*, Z: (-?)(\d{1,9}).*", line)
        if match != None and len(match.groups()) == 6:

            x = -int(match.groups()[1]) if match.groups()[0] == "-" else int(match.groups()[1])
            y = -int(match.groups()[3]) if match.groups()[2] == "-" else int(match.groups()[3])
            z = -int(match.groups()[5]) if match.groups()[4] == "-" else int(match.groups()[5])
            chest_locations.append(ChestLocation(x, y, z))

    if running:
        print("UPDATE seasons SET running = false")
    print(f"INSERT INTO seasons (season_name, running) VALUES ('{season}', '{running}') ON CONFLICT (season_name) DO NOTHING;")
    values = ",".join(f"('{season}', {loc.x}, {loc.y}, {loc.z})" for loc in chest_locations)
    query = f"INSERT INTO chest_locations (season_name, x, y, z) VALUES {values};"
    print(query)

if __name__ == "__main__":
    main()
