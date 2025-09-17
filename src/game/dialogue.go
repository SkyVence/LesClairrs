package game

import (
	"fmt"
	"math"
	"time"
)

// Constante pour la distance de déclenchement du dialogue
const DISTANCE_DECLENCHEMENT = 1.5

// EchangeDialogue représente un tour de la conversation.
type EchangeDialogue struct {
	TextePNJ    string
	TexteJoueur string
}

// Joueur représente le personnage contrôlé par l'utilisateur.
type Joueur struct {
	x float64
	y float64
}

// PNJ représente le personnage non-joueur.
type PNJ struct {
	x float64
	y float64
	// Un drapeau pour indiquer si le dialogue a déjà eu lieu.
	dialoguePasse bool
	// La séquence de dialogue.
	echangeDialogue []EchangeDialogue
}

// NouveauPNJ crée et retourne une nouvelle instance de PNJ.
func NouveauPNJ(x, y float64, echange []EchangeDialogue) *PNJ {
	return &PNJ{
		x:               x,
		y:               y,
		dialoguePasse:   false,
		echangeDialogue: echange,
	}
}

// EstProche vérifie si le joueur est suffisamment proche du PNJ.
func (p *PNJ) EstProche(j *Joueur) bool {
	dx := p.x - j.x
	dy := p.y - j.y
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance <= DISTANCE_DECLENCHEMENT
}

// LancerDialogueTourParTour gère le dialogue tour par tour avec un délai automatique.
func (p *PNJ) LancerDialogueTourParTour(delai time.Duration) {
	fmt.Println("\n--- Le dialogue commence ---")

	// Parcours chaque échange dans le dialogue.
	for _, echange := range p.echangeDialogue {
		// Tour du PNJ
		fmt.Printf("\nPNJ: %s\n", echange.TextePNJ)
		time.Sleep(delai) // Pause après la réplique du PNJ

		// Tour du joueur (réponse prédéfinie)
		fmt.Printf("Vous: %s\n", echange.TexteJoueur)
		time.Sleep(delai) // Pause après la réplique du joueur
	}
	fmt.Println("\n--- Fin du dialogue ---")
	// Marque le dialogue comme terminé.
	p.dialoguePasse = true
}

func main() {
	// Crée les instances du joueur et du PNJ.
	joueur := &Joueur{x: 0, y: 0}
	pnj := NouveauPNJ(5, 0, []EchangeDialogue{
		{
			TextePNJ:    "T’as le coup de main. Qu'est ce que tu fais dans ce trou. T’es sur ma propriété",
			TexteJoueur: "On est où? Vous êtes qui?",
		},
		{
			TextePNJ:    "On est à MEILAND, tu vis dans une cave ? Le meilleur businessman de toute l’histoire Aethelgard en chair et un des derniers complètement en os. Encore une fois... tu fais quoi dans MA mine",
			TexteJoueur: "Je sais pas, je viens de me reveiller et j'ai pas toutes les idées en places",
		},
		{
			TextePNJ:    "T'as le sens de la fête, et le combat dans le fuel. Je te propose un deal, tu jettes un coup d'œil à ma boutique, et je te laisse tranquille. ",
			TexteJoueur: "Merci.",
		},
	})

	fmt.Println("Le joueur se déplace vers le PNJ...")
	fmt.Println("Position initiale du joueur : (0, 0)")
	fmt.Printf("Position du PNJ : (%.1f, %.1f)\n", pnj.x, pnj.y)

	// Définir la durée de l'attente entre chaque tour de dialogue.
	const delaiEntreTours = 2 * time.Second

	// Boucle de jeu principale
	for {
		// Simuler le mouvement du joueur
		if joueur.x < pnj.x {
			joueur.x += 0.5
			fmt.Printf("\nLe joueur avance. Position actuelle : (%.1f, %.1f)\n", joueur.x, joueur.y)
		}

		// Vérifier si le joueur est suffisamment proche et si le dialogue n'a pas encore eu lieu
		if pnj.EstProche(joueur) && !pnj.dialoguePasse {
			fmt.Println("\nLe joueur s'est approché du PNJ.")
			// Lancer le dialogue tour par tour avec le délai spécifié.
			pnj.LancerDialogueTourParTour(delaiEntreTours)
		}

		// Sortir de la boucle une fois le dialogue terminé
		if pnj.dialoguePasse {
			fmt.Println("\nLa conversation est terminée. Le PNJ ne parle plus.")
			break
		}

		// Simuler le temps dans la boucle de jeu
		time.Sleep(1 * time.Second)
	}
}
