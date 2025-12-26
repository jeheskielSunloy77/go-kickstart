import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme } from "@/hooks/use-theme";
import { Computer, Moon, Sun } from "lucide-react";
import type { ReactNode } from "react";

export function ThemeDropdown(props: { children?: ReactNode }) {
  const { theme, setTheme } = useTheme();

  const themeIcons = {
    dark: Moon,
    light: Sun,
    system: Computer,
  };

  const ThemeIcon = themeIcons[theme];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        {props.children || (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="border border-border"
          >
            <ThemeIcon className="size-4" />
          </Button>
        )}
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem onSelect={() => setTheme("system")}>
          <Computer /> System Default
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => setTheme("light")}>
          <Sun /> Light
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => setTheme("dark")}>
          <Moon /> Dark
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
