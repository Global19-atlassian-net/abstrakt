package commands

// visualise is a subcommand that constructs a graph representation of the yaml
// input file and renders this into GraphViz 'dot' notation.
// Initial version renders to dot syntax only, to graphically depict this the output
// has to be run through a graphviz visualisation tool/utiliyy

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/microsoft/abstrakt/internal/dagconfigservice"
	"github.com/microsoft/abstrakt/internal/tools/guid"
)

var verbose string

var visualiseCmd = &cobra.Command{
	Use:   "visualise",
	Short: "format a constellation configuration as Graphviz dot notation",
	Long: "visualise is for producing Graphviz dot notation code of a constellation configuration\n" +
		"abstrakt visualise -f [constellationFilePath]",

	Run: func(cmd *cobra.Command, args []string) {
		if verbose == "true" {
			fmt.Println("args: " + strings.Join(args, " "))
			fmt.Println("constellationFilePath: " + constellationFilePath)
		}

		if !fileExists(constellationFilePath) {
			fmt.Println("Could not open YAML input file for reading")
			os.Exit(-1)
		}

		dsGraph := dagconfigservice.NewDagConfigService()
		err := dsGraph.LoadDagConfigFromFile(constellationFilePath)
		if err != nil {
			log.Fatalf("dagConfigService failed to load file %q: %s", constellationFilePath, err)
		}

		resString := generateGraph(dsGraph)
		fmt.Println(resString)

	},
}

func init() {
	visualiseCmd.Flags().StringVarP(&constellationFilePath, "constellationFilePath", "f", "", "constellation file path")
	err := visualiseCmd.MarkFlagRequired("constellationFilePath")
	if err != nil {
		panic(err)
	}
	visualiseCmd.Flags().StringVarP(&verbose, "verbose", "v", "false", "verbose - show logging  information")
}

//fileExists - basic utility function to check the provided filename can be opened and is not a folder/directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//generateGraph - function to take a dagconfigService structure and create a graph object that contains the
//representation of the graph. Also outputs a string representation (GraphViz dot notation) of the resulting graph
//this can be passed on to GraphViz to graphically render the resulting graph
func generateGraph(readGraph dagconfigservice.DagConfigService) string {

	//lookup is used to map IDs to names. Names are easier to visualise but IDs are more important to ensure the
	//presented constellation is correct and IDs are used to link nodes together
	lookup := make(map[guid.GUID]string)

	g := gographviz.NewGraph()
	if err := g.SetName(strings.Replace(readGraph.Name, " ", "_", -1)); err != nil { //Replace spaces with underscores, names with spaces can break graphviz engines
		log.Fatalf("error: %v", err)
		panic(err)
	}

	if err := g.SetDir(true); err != nil { //Make the graph directed (a constellation is  DAG)
		log.Fatalf("error: %v", err)
		panic(err)
	}

	//Add all nodes to the graph storing the lookup from ID to name (for later adding relationships)
	//Replace spaces in names with underscores, names with spaces can break graphviz engines)
	for _, v := range readGraph.Services {
		if verbose == "true" {
			log.Printf("Adding node %s %s\n", v.ID, v.Name)
		}
		newName := strings.Replace(v.Name, " ", "_", -1)
		lookup[v.ID] = newName
		err := g.AddNode(readGraph.Name, newName, nil)
		if err != nil {
			panic(err)
		}
	}

	//Add relationships to the graph linking using the lookup IDs to name map
	//Replace spaces in names with underscores, names with spaces can break graphviz engines)
	for _, v := range readGraph.Relationships {
		if verbose == "true" {
			log.Printf("Adding relationship from %s ---> %s\n", v.From, v.To)
		}
		localFrom := lookup[v.From]
		localTo := lookup[v.To]
		err := g.AddEdge(localFrom, localTo, true, nil)
		if err != nil {
			panic(err)
		}
	}

	//Produce resulting graph in dot notation format
	// fmt.Printf("%s", g.String())
	return g.String()

}
