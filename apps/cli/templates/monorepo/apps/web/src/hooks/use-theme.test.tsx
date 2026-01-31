import { useTheme } from "@/hooks/use-theme";
import { renderWithCriticalProviders } from "@/testing/render";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { act } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

function ThemeProbe() {
  const { theme, resolvedTheme, setTheme } = useTheme();
  return (
    <div>
      <div data-testid="theme">{theme}</div>
      <div data-testid="resolved">{resolvedTheme}</div>
      <button type="button" onClick={() => setTheme("light")}>
        Light
      </button>
      <button type="button" onClick={() => setTheme("dark")}>
        Dark
      </button>
      <button type="button" onClick={() => setTheme("system")}>
        System
      </button>
    </div>
  );
}

function setupMatchMedia(matches = false) {
  const listeners = new Set<(event: MediaQueryListEvent) => void>();

  const matchMedia = vi.fn().mockImplementation((query: string) => {
    return {
      matches,
      media: query,
      onchange: null,
      addEventListener: (
        _: string,
        cb: (event: MediaQueryListEvent) => void,
      ) => {
        listeners.add(cb);
      },
      removeEventListener: (
        _: string,
        cb: (event: MediaQueryListEvent) => void,
      ) => {
        listeners.delete(cb);
      },
      addListener: () => undefined,
      removeListener: () => undefined,
      dispatchEvent: () => false,
    };
  });

  window.matchMedia = matchMedia as unknown as typeof window.matchMedia;

  return {
    setMatches(nextMatches: boolean) {
      for (const listener of listeners) {
        listener({ matches: nextMatches } as MediaQueryListEvent);
      }
    },
  };
}

beforeEach(() => {
  window.localStorage.clear();
  document.documentElement.classList.remove("dark");
});

describe("useTheme", () => {
  it("defaults to system theme and resolves to system preference", () => {
    setupMatchMedia(true);

    renderWithCriticalProviders(<ThemeProbe />);

    expect(screen.getByTestId("theme").textContent).toBe("system");
    expect(screen.getByTestId("resolved").textContent).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });

  it("persists explicit theme selection", async () => {
    setupMatchMedia(false);

    renderWithCriticalProviders(<ThemeProbe />);

    const user = userEvent.setup();
    await user.click(screen.getByRole("button", { name: "Dark" }));

    expect(screen.getByTestId("theme").textContent).toBe("dark");
    expect(screen.getByTestId("resolved").textContent).toBe("dark");
    expect(window.localStorage.getItem("theme")).toBe("dark");
  });

  it("reacts to system theme changes in system mode", () => {
    const media = setupMatchMedia(false);

    renderWithCriticalProviders(<ThemeProbe />);

    act(() => {
      media.setMatches(true);
    });

    expect(screen.getByTestId("resolved").textContent).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });
});
