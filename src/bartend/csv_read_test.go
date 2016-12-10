package bartend

import (
	"strings"
	"testing"
)

func TestCsvRead(t *testing.T) {
	data :=
		`1,21 - 1, Mixed Cocktail, Tall Highball, Lime Wheel, Mothers Day, 1 1/4 oz. Rum, 1 oz. Lime juice, 1 oz. Simple Syrup, Dash Grenadine, In Shaker, mix ingredients, pour over ice,,,,,
2,21 - 2, Mixed Cocktail, Collins, Orange Slice & Lime Wheel, Green, 1 1/4 oz. Coconut Rum, 3/4 oz. Melon Liqueur, 2 oz. Orange juice, 2 oz. Pineapple juice, Pour Over 1 oz. Soda, In Shaker, mix ingredients, pour over ice,,,,
3,21 - 3, Mixed Cocktail, Highball Glass, Pineapple, Any, 1 1/2 oz. Early Times, 1 oz. Apricot Brandy, 2 oz. Pineapple juice, Club Soda, Shake ingredients with ice; pour into highball glass; stir in club soda,,,,,,,
4,57 T - Birds (California Plates), Mixed Cocktail, Collins, None, Any, 1 1/4 oz. Scotch, 1 oz. Coffee Liqueur, 3 oz. Half & Half, In Shaker, mix ingredients, pour over ice,,,,,,
5,57 T - Birds (Cape Cod Plates), Shooter, Chilled Rocks, None, Any, 1/2 oz. Orange Brandy Liqueur, 1/2 oz.Amaretto, 1/2 oz. Vodka 1/2 oz. Cranberry juice, In Shaker, mix ingredients with ice and strain into glass,,,,,,,
6,57 T - Birds (Florida Plates), Shooter, Chilled Rocks, None, Any, 1/2 oz. Orange Brandy Liqueur, 1/2 oz.Amaretto, 1/2 oz. Vodka 1/2 oz. Grapefruit juice, In Shaker, mix ingredients with ice and strain into glass,,,,,,,
7,57 T - Bird, Mixed Cocktail, Rocks Glass, None, Any, 1/2 oz. Southern Comfort, 1/2 oz. Grand Marnier, 1/2 oz. Amaretto Di Saronno, splash Pineapple juice, Fill shaker half full with ice, add ingredients; stir; strain into rocks glass,,,,,,
8,77 Sunset Strip - 1, Martini, Chilled Cocktail, 3 Speared Cocktail Onions, Fathers Day, 2 oz. Gin or Vodka, Dash Dry Vermouth, In Shaker, mix ingredients with ice and strain into glass,,,,,,,,
9,77 Sunset Strip - 2, Shooter, Chilled Rocks, None, Any, 1/2 oz. Orange Brandy Liqueur, 1/2 oz. Amaretto, 1/2 oz. Vodka, 1/2 oz. Pineapple juice, In Shaker, mix ingredients with ice and strain into glass,,,,,,
10,77 Sunset Strip - 3, Shooter, Cordial, None, 4th of July, 1/2 oz. Creme de Noya, 1/2 oz. White Creme de Menthe, 1/2 oz. Blue Curacao, Pour Liqueurs in the order they are listed,,,,,,,,
11,7th - Inning Stretch, Mixed Cocktail, Cocktail Glass, None, Any, 1 part Jack Daniels, 2 parts Orange juice, 2 parts 7-Up, Mix ingredients and pour into cocktail glass,,,,,,,,
12,Adams Apple - 1, Highball, Tall Rocks, None, Any, 1 1/2 oz. Galliano, 3/4 oz. White Creme de Menthe, Pour ingredients as listed over ice,,,,,,,,,
13,Adams Apple - 2, Martini, Chilled Cocktail, Lemon Twist, Any, 1 oz. Harveys Bristol Cream, 1 oz. Dubonnet, In Shaker, mix ingredients with ice and strain into glass,,,,,,,,
14,Adams Apple - 3, Mixed Cocktail, Rocks Glass, None, Any, 1 1/4 oz. Finlandia Vodka, Pour Finlandia vodka over ice, fill with apple cider, stir,,,,,,,,
15,Adios Mother - 1, Frozen Drink, Tall Specialty, Lime Wheel, Cinco de Mayo, 1 1/4 oz. Gold Tequila, 1/2 oz. Triple Sec, 1 1/2 oz. Margarita Mix, 1 Scoop Crushed Ice, Combine ingredients in blender, blend until smooth,,,,,,
16,Adios Mother - 2, Highball, Tall Highball, None, Any, 3/4 oz. Irish Cream, 3/4 oz. Amaretto, 1 oz. Half & Half, In Shaker, mix ingredients, pour over ice,,,,,,
17,Adios Mother - 3, Mixed Cocktail, Cocktail Glass, None, Mothers Day, 1/2 oz. Finlandia Vodka, 1/2 oz. Gin, 1/2 oz. Rum, 3 oz. Sweet and Sour Mix, Blend for a few seconds,,,,,,,
18,Adios Mother - 4, Mixed Cocktail, Collins, None, Halloween, 1 1/4 oz. Vodka, 4 oz. Apple Cider, Pour ingredients as listed over ice,,,,,,,,,
19,Admiral Jack, Mixed Cocktail, Cocktail Glass, None, Any, One part Jack Daniels, 1/2 part Noilly Prat Dry Vermouth, Lemon juice, Mix ingredients as listed,,,,,,,,
20,After 5, Shooter, Cordial, None, Any, 1/2 oz. Coffee Liqueur, 1/2 oz. Peppermint Schnapps, 1/2 oz. Irish Cream, Pour Liqueurs in the order they are listed,,,,,,,,
21,Airport - 1, Mixed Cocktail, Chimney Glass, Orange slice, Any, 1 1/2 oz. Jack Daniels, 1/2 oz. Orange juice, dash Lemon juice, Asti Spumante Wine, Combine Jack Daniels and orange juice in glass; top with Asti Spumante,,,,,,,
22,Airport - 2, Shooter, Chilled Rocks, None, Any, 1 oz. Vodka, 1/2 oz. Triple Sec, juice of 1/2 Lemon, 1/2 oz. Simple Syrup, In Shaker, mix ingredients with ice and strain into glass,,,,,,
23,Alabama Slammer - 1, Mixed Cocktail, Collins, None, Valentines Day, 3/4 oz. Vodka, 3/4 oz. Bourbon, 3/4 oz. Amaretto, 3 oz. Orange juice, 1/2 oz. Grenadine, In Shaker, mix ingredients, pour over ice,,,,
24,Alabama Slammer - 2, Shooter, Shot Glass, None, Any, 1 part Southern Comfort, 1 part Club Soda, Place napkin over glass and slam glass on bar top,,,,,,,,,
25,Alabama Slammer - 3, Shooter, Shot Glass, None, Any, 3/4 oz. Southern Comfort, 1 oz. Pepe Lopez Tequila, 1/2 oz. Orange juice, 1/2 oz. Cranberry juice, All drinks must be shaken with rocks and strained into shot glass,,,,,,,
26,Alabama Slammer - 4, Martini, Chilled Cocktail, Lime Wheel, Cinco de Mayo, 1 1/4 oz. Tequila, 1/2 oz. Triple Sec, 1 oz. Margarita Mix, 1 oz. Water, Salt Rim, In Shaker, mix ingredients with ice and strain into glass,,,,,
27,Alabama Slammer - 5, Mixed Cocktail, Collins, None, Luau, 1/2 oz. Gold Rum, 1/2 oz. Coconut Rum, 1/2 oz. Banana Liqueur, 1/2 oz. Grenadine, 1 oz. Orange juice, 1 oz. Pineapple juice, In Shaker, mix ingredients, pour over ice,,,
28,Alabama Slammer - 6, Punch, Tall Glass, Orange slice, Birthday, 1 oz. Southern Comfort, 1/2 oz. Sloe gin, 1/2 oz. Amaretto, 1 1/2 oz. Orange juice, Fill tall glass with ice and add all ingredients,,,,,,,
29,Alejandros - 1, Hot Drink, Footed Glass Mug, Preheated, Heavy Whipping Cream, Winter, 1 1/4 oz. Dark Rum, Fill with Hot Tea, Pour ingredients as listed into preheated mug,,,,,,,,
30,Alejandros - 2, Mixed Cocktail, Cocktail Glass, Lemon and cherry, Any, 1 1/2 oz. Southern Comfort, 2 oz. Pineapple juice, 2 oz. Orange juice, 1 oz. Cranberry juice, splash Club Soda, Speed shake,,,,,,
31,Alejandros - 3, Mixed Cocktail, Collins, Lemon Squeeze, 4th of July, 3/4 oz. Bourbon, 3/4 oz. Triple Sec, 3/4 oz. Sweet & Sour Mix, Fill with Lemon-Lime Soda, In Shaker, mix ingredients, pour over ice,,,,,
32,Algonquin, Mixed Cocktail, Cocktail Glass, None, Any, 1 1/2 oz. Early Times, 1 oz. Noilly Prat Dry Vermouth, 1 oz. Pineapple juice, Mix ingredients as listed,,,,,,,,
33,Alice In Wonderland - 1, Shooter, Chilled Rocks, None, Any, 1/2 oz. Coffee Liqueur, 1/2 oz. Tequila, 1/2 oz. Orange Brandy Liqueur, In Shaker, mix ingredients with ice and strain into glass,,,,,,,
34,Alice In Wonderland - 2, Shooter, Shot Glass, None, Any, 1 part Pepe Lopez Tequila, 1 part Grand Marnier, 1 part Tia Maria, Serve as shot,,,,,,,,
35,Almond Cappuccino, Cappuccino, Footed Glass Mug, Heavy Whipping Cream, Shaved Chocolate & Cookie, Any, 1 oz. Amaretto, 3 oz. Milk, 1 Espresso Shot, Steam milk and liquor together in mug until foamy, Brew 1 Shot of Espresso, Pour Shot of Espresso into center of steamed milk/liquor mixture,,,,,
36,Almond Joey, Ice Cream Drink, Tall Specialty, None, Any, 1 1/4 oz. Amaretto, 1 oz. Cream of Coconut, 1 oz. Chocolate Syrup, 2 Scoops Vanilla Ice Cream, 1/2 Scoop Crushed Ice, Combine ingredients in blender, blend until smooth,,,,,
37,Almond Lemonade, Mixed Cocktail, Tall Glass, Lemon slice, Any, 1 1/4 oz. Finlandia Vodka, 1/4 oz. Creme de Almond, Lemonade, Pour over ice in a tall class,,,,,,,,
38,Almond Mocha, Coffee Drink, Footed Glass Mug, Preheated, Heavy Whipping Cream, Any, 1 1/4 oz. Amaretto, Fill with Hot Coffee, Pour ingredients as listed into preheated mug,,,,,,,,
39,Almond Tea, Hot Drink, Footed Glass Mug, Preheated, None, Any, 1 Orange Spice Tea Bag, Fill with Hot Water and Brew, 1 1/4 oz. Amaretto, Pour ingredients as listed into preheated mug. Leave bag in glass,,,,,,,
40,Als Special Cocktail, Mixed Cocktail, Cocktail Glass, None, Any, 2/3 oz. Canadian Mist, 1/3 oz. Apricot Brandy, Mix ingredients as listed,,,,,,,,,
41,Amaretto Colada, Frozen Drink, Tall Specialty, Pineapple Slice & Cherry, Any, 1 1/4 oz. Amaretto, 5 oz. Pina Colada Mix, 1 Scoop Crushed Ice, Combine ingredients in blender, blend until smooth,,,,,,,`
	rdr := NewReader(strings.NewReader(data))

	for {
		r, err := rdr.Read()
		if r == nil {
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		t.Log(r)
	}
}
