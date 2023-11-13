package page

// WeightedPages is a list of Pages with their corresponding (and relative) weight
// [{Weight: 30, Page: *1}, {Weight: 40, Page: *2}]
type WeightedPages []WeightedPage

type WeightedPage struct {
  Weight int
  Page

  // Referenct to the owning Page. This avoids having to do
  // manual. Site.GetPage lookups. It is implemented in this roundabout way
  // because we cannot add additional state to the WeightedPages slice
  // without breaking lots of templates in the wild.
  owner Page
}
