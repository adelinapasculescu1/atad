package cli

import (
    "fmt"

    "github.com/spf13/cobra"
)

func NewCategoryCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "category",
        Short: "Manage categories",
    }

    cmd.AddCommand(NewCategoryAddCommand(deps))
    cmd.AddCommand(NewCategoryListCommand(deps))
    return cmd
}

func NewCategoryAddCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "add",
        Short: "Add a new category",
        RunE: func(cmd *cobra.Command, args []string) error {
            name, _ := cmd.Flags().GetString("name")
            if name == "" {
                return fmt.Errorf("--name is required")
            }

            c, err := deps.CategoryRepo.Create(name)
            if err != nil {
                return err
            }

            fmt.Printf("Category created: %d - %s\n", c.ID, c.Name)
            return nil
        },
    }

    cmd.Flags().String("name", "", "Category name")
    return cmd
}

func NewCategoryListCommand(deps Deps) *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "List categories",
        RunE: func(cmd *cobra.Command, args []string) error {
            cats, err := deps.CategoryRepo.List()
            if err != nil {
                return err
            }

            if len(cats) == 0 {
                fmt.Println("No categories found.")
                return nil
            }

            for _, c := range cats {
                fmt.Printf("%d\t%s\n", c.ID, c.Name)
            }
            return nil
        },
    }
}
