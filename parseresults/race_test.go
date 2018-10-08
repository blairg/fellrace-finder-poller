package parseresults

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestParseRace5911(t *testing.T) {
	content, err := ioutil.ReadFile("../testResources/race-5911.html")
	raceHTML := string(content)

	//fmt.Println(raceHTML)

	if err != nil {
		log.Fatal(err)
	}

	race := ParseRace("5911", raceHTML)
	jsonRace, err := json.Marshal(race)

	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(string(jsonRace))

	// jsonRecipeString := cleanUpJSON(string(jsonRecipe))
	// expectedValue := "nameAngelasslowroastedgingerporkdescriptionpreparationTimelessthan30minscookingTimeover2hoursservesServes68ingredientsNameIngredient1shoulderofpork4garlic4cm1½inchpiecefreshrootginger3tbspoliveoil4tbspwhitewinevinegarmethodPreheattheovento220C425FGasmark7PlacetheporkskinsideuponarackoveraroastingtinPlacethegarlicandgingerinapestleandmortarorfoodprocessorandpoundorprocessuntilyougetaroughpastethenmixintheoilandvinegarRubthepastealloverthescoredskinoftheporkPlaceinthepreheatedovenandcookfor30minutesRemovetheporkfromtheovenreducethetemperatureto150C300FGas2Turntheporkoverwiththeskinsidedownontherackandreturntotheovenandcookfor45hoursRemovefromtheovenandturnuptothehighestsetting220C425FGas7Turntheporkovertothecracklingsideontherackandroastinthehotovenforthefinal20minutestocrispupthecracklingLeavetostandfor1015minutesbeforecarvingToservecutawaythecracklingwithasharpknifeandbreakitupintopiecesthencarvethemeatItshouldbeverytenderandsucculentchefAngelaBoggianourlhttpswwwbbccomfoodrecipesangelasslowroastedgi_71381ampUrlhttpswwwbbccomfoodrecipesangelasslowroastedgi_71381ampimageUrl"

	// if jsonRecipeString != expectedValue {
	// 	t.Errorf("Expected value does not match")
	// }
}
