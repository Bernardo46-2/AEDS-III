# AEDS-III
Repositório para os trabalhos feitos em Algoritmos e Estruturas de Dados III

O programa segue a seguinte hierarquia de metodos:
main >> middlewares >> handlers >> crud >> dataManager >> models/data

Middlewares = Interceptam as requisiçoes e fazem o tratamento de cors
Handlers = Fazem o tratamento http e chamam as respectivas funções de crud e ordenação
Crud = Realiza a conversa entre o DataManager e o Models
Data Manager = Controle do arquivo binario, de crud a ordenação
Models = Structs e funções de controle da struct
Data = Arquivo binario e CSV original
Logger = Log do servidor
Deprecated = Codigos nao mais utilizados