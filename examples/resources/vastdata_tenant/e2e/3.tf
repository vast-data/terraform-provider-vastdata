resource "vastdata_tenant" "with_views_count" {
  name            = "vastdbtenant-views"
  get_views_count = true
}
