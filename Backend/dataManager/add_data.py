import csv
import urllib.request

IN_FILE = '../data/pokedex.csv'
OUT_FILE = '../data/pokedex2.csv'
URL = 'https://raw.githubusercontent.com/PokeAPI/pokeapi/master/data/v2/csv/pokemon_species_flavor_text.csv'

def get_descriptions() -> str:
    try:
        with urllib.request.urlopen(URL) as f:
            return f.read().decode('utf-8')
    except urllib.error.URLError as e:
        return e.reason
    

def parse_description(string: str, id: int, start: int) -> tuple[int, str]:
    token = f'{id},18,9'
    index = string.find(token, start) + len(token) + 2
    description = ''

    while string[index] != '"':
        if string[index] != '\n':
            description += string[index]
        else:
            description += ' '
        index += 1

    return index, description


if __name__ == '__main__':
    with open(OUT_FILE, 'w') as new_file:
        file = open(IN_FILE)
        csv_file = file.readlines()
        file.seek(0)
        csv_reader = csv.DictReader(file)
        descriptions = get_descriptions()
        last_index = 0
        last_id = 0
        i = 0

        headers = csv_file[i]
        new_file.write(headers[:len(headers) - 1] + ',description\n')

        for row in csv_reader:
            i += 1
            poke_id = row.get('pokedex_number')

            if poke_id == None or poke_id == last_id:
                continue

            last_id = poke_id
            last_index, description = parse_description(descriptions, poke_id, last_index)

            line = csv_file[i]
            new_file.write(line[:len(line) - 1] + f',"{description}"\n')

        file.close()
