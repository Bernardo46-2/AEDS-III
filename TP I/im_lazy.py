import csv
import requests
import json

URL = 'https://pokeapi.glitch.me/v1/pokemon/'
GEN_RELEASE_DATES = ['February 27, 1996', 
                     'November 21, 1999', 
                     'November 21, 2002',
                     'September 28, 2006',
                     'September 18, 2010',
                     'October 12, 2013',
                     'November 18, 2016',
                     'November 15, 2019',
                     'November 18, 2022',
                     'October 24, 1929']


def make_req(pokemon):
    req = requests.get(URL + pokemon).text
    gen = json.loads(req)

    if isinstance(gen, list):
        gen = gen[0]

    return gen.get('gen')


if __name__ == '__main__':
    with open('csv/edited_pokedex.csv', 'w') as new_file:
        file = open('csv/pokedex.csv')
        csv_file = file.readlines()
        file.seek(0)
        csv_reader = csv.DictReader(file)
        i = 0

        headers = csv_file[i]
        new_file.write(headers[:len(headers) - 1] + ',gen,release_date\n')

        for row in csv_reader:
            i += 1
            pokemon = row.get('name')

            if pokemon == None:
                continue

            gen = make_req(pokemon)

            if gen == None:
                gen = 10

            line = csv_file[i]
            new_file.write(line[:len(line) - 1] + f',{gen},"{GEN_RELEASE_DATES[int(gen) - 1]}"\n')

        file.close()
