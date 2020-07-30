# Le Compte Est Bon

Le jeu du **Compte est Bon** est un jeu basé sur le calcul mental, issu du jeu télévisé ** Des Chiffres et des Lettres**. A partir d'un tirage de 6 plaques, il faut trouver un nombre entre 101 et 999, en utilisant les 4 opérations arithmétiques d'addition, soustraction, multiplication ou division. Les résultats intermédiaires doivent être des entiers positifs.

Cette première version est réalisée en Java et me permet d'en pratiquer un peu.  
L'algorithme est un *brute force*: le programme essaie toutes les combinaisons de calcul possible, entre toutes les 6 plaques du tirage initial, ce qui donne ensuite un jeu de 5 nombres, pour chacune des 4 opérations, puis le programme itère récursivement jusqu'à trouver le compte exact. La meilleure solution approchée est gardée à chaque étape, pour le cas où le compte exact ne peut pas être trouvé.  
Quelques optimisations permettent de réduire le nombre de combinaisons possibles:
- comme les soustractions et les divisions doivent produire des entiers positifs, le programme classe les nombres dans l'ordre croissant
- les additions et les multiplications sont commutatives, il suffit de faire **A + B**, pas besoin de refaire **B + A**
- les soustractions qui aboutissent à un résultat **inférieur ou égal à 0** sont inutiles
- les multiplications et les divisions par **1** sont inutiles
- les divisions avec un résultat non entier sont inutiles

