package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/iwittkau/projects"
	input "github.com/tcnksm/go-input"

	"github.com/iwittkau/projects/configuration"
	"github.com/iwittkau/projects/project"

	. "github.com/logrusorgru/aurora"

	flag "github.com/ogier/pflag"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var (
		add bool
		git bool
	)

	flag.BoolVar(&add, "add", false, "add current folder to .projects file")
	flag.BoolVar(&git, "git", false, "analyze git repo commits, falls back to file attributes for non-repository folders")
	flag.Parse()

	p, err := configuration.ReadConfigFromHomeDir()
	if err != nil {
		panic(err)
	}
	if add {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		query := "Path to project?"
		path, err := ui.Ask(query, &input.Options{
			Default:  wd,
			Required: true,
			Loop:     true,
		})
		query = "Project name?"
		name, err := ui.Ask(query, &input.Options{
			Default:  filepath.Base(wd),
			Required: true,
			Loop:     true,
		})

		ws := projects.Workspace{
			Path:   path,
			Name:   name,
			Active: true,
		}

		if err := configuration.AddProject(&p, ws); err != nil {
			fmt.Println(Red(err.Error()))
			os.Exit(1)
		}

		if err := configuration.WriteToHomeDir(p); err != nil {
			fmt.Println(Red(err.Error()))
			os.Exit(1)
		}

		fmt.Println(Sprintf(Green("Project '%s' added to '~/.projects'"), ws.Name))

		os.Exit(0)
	}

	sort.Sort(projects.ByWorkSpaceName(p.Workspaces))

	fmt.Print("\nProjects:\n\n")

	var (
		totalSize uint64
		oldest    = time.Now()
		oldestI   int
		latest    time.Time
		latestI   int
		inactive int
		pros      []projects.Project
	)
	for i := range p.Workspaces {
		pro, err := project.FromWorkspace(p.Workspaces[i])
		if err != nil {
			log.Println(err.Error())
		} else {
			pros = append(pros, pro)
		}
	}

	activity := map[projects.Project]string{}

	for i := range pros {
		if !pros[i].Active{
			inactive++
		}
		totalSize = totalSize + uint64(pros[i].SizeBytes)
		var dates []time.Time
		var gitDates []time.Time
		var err error
		if git {
			dates, err = project.ListDates(pros[i].Path)
			if err != nil {
				log.Println(err)
			}
			gitDates, _ = project.ListCommits(pros[i].Path)

		} else {
			dates, err = project.ListDates(pros[i].Path)
			if err != nil {
				log.Println(err)
			}
		}

		datesStr, latest, err := analyzeActivity(dates, 'o', nil, 18)
		if err != nil {
			log.Println(err.Error())
		}
		if git {
			// log.Println(len(dates), len(gitDates))
			var lt time.Time
			datesStr, lt, err = analyzeActivity(gitDates, '+', datesStr, 18)
			if err != nil {
				log.Println(err.Error())
			}
			if lt.After(latest) {
				latest = lt
			}
		}

		if pros[i].LastEdited.Before(latest) {
			pros[i].LastEdited = latest
		}

		activity[pros[i]] = string(datesStr)
	}

	sort.Sort(projects.ByProjectLastEdited(pros))

	dateF := "2006-01-02"

	for i := range pros {

		iStr := ""
		if i < 9 {
			iStr = fmt.Sprintf("0%d", i+1)
		} else {
			iStr = fmt.Sprintf("%d", i+1)
		}

		proStr := fmt.Sprintf("%s [%s]  %s  %s", iStr, activity[pros[i]], pros[i].LastEdited.Format(dateF), pros[i].Name)

		if pros[i].Active {

			if oldest.After(pros[i].LastEdited) {
				oldest = pros[i].LastEdited
				oldestI = i
			}
			if latest.Before(pros[i].LastEdited) {
				latest = pros[i].LastEdited
				latestI = i
			}

			switch {
			case time.Since(pros[i].LastEdited) < time.Hour*24*30:
				fmt.Println(Green(proStr))
			case time.Since(pros[i].LastEdited) > time.Hour*24*30 && time.Since(pros[i].LastEdited) < time.Hour*24*90:
				fmt.Println(Cyan(proStr))
			case time.Since(pros[i].LastEdited) > time.Hour*24*90 && time.Since(pros[i].LastEdited) < time.Hour*24*180:
				fmt.Println(Brown(proStr))
			case time.Since(pros[i].LastEdited) > time.Hour*24*180 && time.Since(pros[i].LastEdited) < time.Hour*24*365:
				fmt.Println(Red(proStr))
			case time.Since(pros[i].LastEdited) > time.Hour*24*365:
				fmt.Println(Magenta(proStr))
			default:
			}
		} else {
			fmt.Println(Gray(125, proStr))
		}

	}

	fmt.Printf(Sprintf("\n%s %s %s %s %s %s\n\n",
		Inverse(Green(" < 30 days ")),
		Inverse(Cyan(" > 30 days ")),
		Inverse(Brown(" > 90 days ")),
		Inverse(Red(" > 180 days ")),
		Inverse(Magenta(" > 365 days ")),
		Inverse(Gray(125, " Inactive "))))
	if len(pros) < 1 {
		os.Exit(0)
	}
	fmt.Println("Oldest:", fmt.Sprintf("%d days", int(time.Since(oldest).Truncate(time.Hour).Hours()/24.0)), fmt.Sprintf("(%s)", pros[oldestI].Name))
	fmt.Println("Latest:", time.Since(latest).Truncate(time.Hour).Hours(), "hours", fmt.Sprintf("(%s)", pros[latestI].Name))
	fmt.Println("Active:", len(pros)-inactive)
	fmt.Println("Inactive:", inactive)

	// err = configuration.WriteToFile(p)
	// if err != nil {
	// 	panic(err)
	// }

}

func analyzeActivity(dates []time.Time, rn rune, existing []rune, months int) ([]rune, time.Time, error) {
	moreRn := '<'
	fillRn := '.'
	spaceRn := ' '
	var existingRn rune
	activity := make([]rune, months)
	if existing != nil {
		if len(existing) != months {
			return nil, time.Time{}, fmt.Errorf("length of existing (%d) and amount of months (%d) do not match", len(existing), months)
		}
		activity = existing
		for i := range activity {
			if activity[i] != moreRn && activity[i] != fillRn && activity[i] != spaceRn {
				existingRn = activity[i]
			}
		}
	} else {

		for i := 0; i < len(activity); i++ {
			activity[i] = spaceRn
		}
	}
	latest := time.Time{}

	for i := range dates {
		if latest.Before(dates[i]) {
			latest = dates[i]
		}
		monthsAgo := time.Since(dates[i]) / (time.Hour * 24 * (365 / 12))
		month := int(monthsAgo)
		if month < months {
			activity[months-1-month] = rn
		} else {
			activity[0] = moreRn
		}
	}

	first := false
	for i := 0; i < len(activity); i++ {
		if activity[i] != ' ' {
			first = true
		}
		if first && activity[i] != rn && activity[i] != moreRn && activity[i] != existingRn {
			activity[i] = fillRn
		}
	}

	return activity, latest, nil
}
