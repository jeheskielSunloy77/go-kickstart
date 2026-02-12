package ui

import (
	"fmt"
	"strings"
)

const welcomeASCII = `
 ██████╗  ██████╗       ██╗  ██╗██╗ ██████╗██╗  ██╗███████╗████████╗ █████╗ ██████╗ ████████╗
██╔════╝ ██╔═══██╗      ██║ ██╔╝██║██╔════╝██║ ██╔╝██╔════╝╚══██╔══╝██╔══██╗██╔══██╗╚══██╔══╝
██║  ███╗██║   ██║█████╗█████╔╝ ██║██║     █████╔╝ ███████╗   ██║   ███████║██████╔╝   ██║
██║   ██║██║   ██║╚════╝██╔═██╗ ██║██║     ██╔═██╗ ╚════██║   ██║   ██╔══██║██╔══██╗   ██║
╚██████╔╝╚██████╔╝      ██║  ██╗██║╚██████╗██║  ██╗███████║   ██║   ██║  ██║██║  ██║   ██║
 ╚═════╝  ╚═════╝       ╚═╝  ╚═╝╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝
`

func WelcomeCard(version string, repoURL string) string {
	lines := []string{
		BannerStyle().Render(strings.TrimRight(welcomeASCII, "\n")),
		SectionTitleStyle().Render("⚙️  go-kickstart // meme arcade mode engaged"),
		HintStyle().Render(fmt.Sprintf("Version: %s", version)),
		HintStyle().Render(fmt.Sprintf("Repo: %s", repoURL)),
		HintStyle().Render(ContributionLine()),
	}

	return strings.Join(lines, "\n")
}
