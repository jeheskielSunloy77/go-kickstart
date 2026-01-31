import { useTheme } from "@/hooks/use-theme";
import {
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@go-kickstart/ui";
import { Computer, Moon, Sun } from "lucide-react";
import type { ReactNode } from "react";

export function ThemeDropdown(props: { children?: ReactNode }) {
  const { setTheme, resolvedTheme } = useTheme();

  const themeIcons = {
    dark: Moon,
    light: Sun,
  };

  const ThemeIcon = themeIcons[resolvedTheme];

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
