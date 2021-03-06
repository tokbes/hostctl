package file

import (
	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/render"
)

// ListOptions contains available options for listing.
type ListOptions struct {
	Renderer     render.Renderer
	Profiles     []string
	ProfilesOnly bool
	StatusFilter host.Status
}

func includeProfile(needle string, stack []string) bool {
	if len(stack) == 0 {
		return true
	}

	for _, s := range stack {
		if s == needle {
			return true
		}
	}

	return false
}

// ProfileStatus shows a table only with profile names status
func (f *File) ProfileStatus(r render.Renderer, profiles []string) {
	for _, name := range f.data.ProfileNames {
		currProfile := f.data.Profiles[name]

		if profiles != nil && !includeProfile(name, profiles) {
			continue
		}

		r.AppendRow(&render.Row{
			Profile: currProfile.Name,
			Status:  currProfile.GetStatus(),
		})
	}

	_ = r.Render()
}

// List shows a table with profile names status and routing information
func (f *File) List(r render.Renderer, opts *ListOptions) {
	addDefault(f, r, opts)

	for _, name := range f.data.ProfileNames {
		addProfiles(f.data.Profiles[name], r, opts)
	}

	_ = r.Render()
}

func addDefault(f *File, r render.Renderer, opts *ListOptions) {
	// First check if default should be shown
	if !includeProfile(host.Default, opts.Profiles) {
		return
	}

	for _, row := range f.data.DefaultProfile {
		if row.Comment == "" && row.Profile != "" {
			r.AppendRow(row)
		}
	}

	if len(f.data.DefaultProfile) > 0 && len(f.data.Profiles) > 0 {
		r.AddSeparator()
	}
}

func addProfiles(p *host.Profile, r render.Renderer, opts *ListOptions) {
	if !includeProfile(p.Name, opts.Profiles) {
		return
	}

	if opts.StatusFilter != "" && p.Status != opts.StatusFilter {
		return
	}

	for _, ip := range p.IPList {
		route := p.Routes[ip]
		for _, h := range route.HostNames {
			r.AppendRow(&render.Row{
				Profile: p.Name,
				Status:  p.GetStatus(),
				IP:      route.IP.String(),
				Host:    h,
			})
		}
	}

	if len(p.IPList) > 0 {
		r.AddSeparator()
	}
}
