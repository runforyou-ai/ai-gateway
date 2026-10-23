import type { ReactNode } from "react";
import {
  useGetLocale,
  useLink,
  useLogout,
  useRefineOptions,
  useSetLocale,
  useTranslation,
} from "@refinedev/core";
import { UserAvatar } from "@/components/refine-ui/layout/user-avatar";
import { useTheme } from "@/components/refine-ui/theme/theme-provider";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useSidebar, SidebarTrigger } from "@/components/ui/sidebar";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { Check, CircleUser, Languages, LogOut, Monitor, Moon, Sun } from "lucide-react";

export const Header = () => {
  const { isMobile } = useSidebar();

  return <>{isMobile ? <MobileHeader /> : <DesktopHeader />}</>;
};

function DesktopHeader() {
  const { open } = useSidebar();

  return (
    <header
      className={cn(
        "sticky",
        "top-0",
        "flex",
        "h-16",
        "shrink-0",
        "items-center",
        "gap-4",
        "border-b",
        "border-border",
        "bg-sidebar",
        "px-3",
        "z-40"
      )}
    >
      {!open ? (
        <SidebarTrigger className={cn("text-muted-foreground")} />
      ) : null}
      <div className={cn("ml-auto", "flex", "items-center", "gap-2")}>
        <LanguageSwitcher />
        <UserAccountMenu />
      </div>
    </header>
  );
}

function MobileHeader() {
  const { open, isMobile } = useSidebar();

  const { title } = useRefineOptions();

  return (
    <header
      className={cn(
        "sticky",
        "top-0",
        "flex",
        "h-12",
        "shrink-0",
        "items-center",
        "gap-2",
        "border-b",
        "border-border",
        "bg-sidebar",
        "pr-3",
        "justify-between",
        "z-40"
      )}
    >
      <SidebarTrigger
        className={cn("text-muted-foreground", "rotate-180", "ml-1", {
          "opacity-0": open,
          "opacity-100": !open || isMobile,
          "pointer-events-auto": !open || isMobile,
          "pointer-events-none": open && !isMobile,
        })}
      />

      <div
        className={cn(
          "whitespace-nowrap",
          "flex",
          "flex-row",
          "h-full",
          "items-center",
          "justify-start",
          "gap-2",
          "transition-discrete",
          "duration-200",
          {
            "pl-3": !open,
            "pl-5": open,
          }
        )}
      >
        <div>{title.icon}</div>
        <h2
          className={cn(
            "text-sm",
            "font-bold",
            "transition-opacity",
            "duration-200",
            {
              "opacity-0": !open,
              "opacity-100": open,
            }
          )}
        >
          {title.text}
        </h2>
      </div>

      <div className={cn("flex", "items-center", "gap-1")}>
        <LanguageSwitcher triggerClassName={cn("h-8", "w-8")} />
        <UserAccountMenu triggerClassName={cn("h-8", "w-8")} />
      </div>
    </header>
  );
}

type HeaderIconButtonProps = {
  triggerClassName?: string;
};

function LanguageSwitcher({ triggerClassName }: HeaderIconButtonProps) {
  const { translate } = useTranslation();
  const getLocale = useGetLocale();
  const setLocale = useSetLocale();
  const current = (getLocale() ?? "en").toLowerCase().startsWith("zh")
    ? "zh"
    : "en";

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          aria-label={translate("language.switch")}
          className={cn(
            "cursor-pointer",
            "select-none",
            "rounded-full",
            "text-muted-foreground",
            "hover:bg-transparent",
            "hover:text-foreground",
            "dark:hover:bg-transparent",
            "focus-visible:ring-0",
            "focus-visible:border-transparent",
            "data-[state=open]:bg-transparent",
            "data-[state=open]:cursor-pointer",
            triggerClassName
          )}
        >
          <Languages className={cn("h-[1.2rem]", "w-[1.2rem]")} />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className={cn("min-w-36", "rounded-lg")}>
        <DropdownMenuItem
          className={cn("cursor-pointer", "justify-between")}
          onClick={() => setLocale("zh")}
        >
          {translate("language.zh")}
          {current === "zh" ? <Check className={cn("h-4", "w-4")} /> : null}
        </DropdownMenuItem>
        <DropdownMenuItem
          className={cn("cursor-pointer", "justify-between")}
          onClick={() => setLocale("en")}
        >
          {translate("language.en")}
          {current === "en" ? <Check className={cn("h-4", "w-4")} /> : null}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function UserAccountMenu({ triggerClassName }: HeaderIconButtonProps) {
  const { theme, setTheme } = useTheme();
  const { mutate: logout, isPending: isLoggingOut } = useLogout();
  const { translate } = useTranslation();
  const Link = useLink();

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className={cn(
            "cursor-pointer",
            "select-none",
            "rounded-full",
            "p-0",
            "hover:bg-transparent",
            "dark:hover:bg-transparent",
            "focus-visible:ring-0",
            "focus-visible:border-transparent",
            "data-[state=open]:bg-transparent",
            "data-[state=open]:cursor-pointer",
            triggerClassName
          )}
        >
          <UserAvatar />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className={cn("w-48", "rounded-lg", "p-1")}>
        <DropdownMenuItem asChild>
          <Link to="/console/profile" className={cn("cursor-pointer")}>
            <CircleUser />
            {translate("account.profile")}
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <div className={cn("flex", "items-center", "justify-between", "gap-1", "px-1", "py-1.5")}>
          <ThemeIconButton
            active={theme === "light"}
            label={translate("theme.switchToLight")}
            onClick={() => setTheme("light")}
          >
            <Sun className={cn("h-4", "w-4")} />
          </ThemeIconButton>
          <ThemeIconButton
            active={theme === "dark"}
            label={translate("theme.switchToDark")}
            onClick={() => setTheme("dark")}
          >
            <Moon className={cn("h-4", "w-4")} />
          </ThemeIconButton>
          <ThemeIconButton
            active={theme === "system"}
            label={translate("theme.switchToSystem")}
            onClick={() => setTheme("system")}
          >
            <Monitor className={cn("h-4", "w-4")} />
          </ThemeIconButton>
        </div>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          disabled={isLoggingOut}
          onClick={() => logout()}
        >
          <LogOut />
          {isLoggingOut
            ? translate("buttons.loggingOut")
            : translate("buttons.logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

type ThemeIconButtonProps = {
  active: boolean;
  label: string;
  onClick: () => void;
  children: ReactNode;
};

function ThemeIconButton({
  active,
  label,
  onClick,
  children,
}: ThemeIconButtonProps) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <button
          type="button"
          aria-label={label}
          aria-pressed={active}
          onClick={onClick}
          onPointerDown={(event) => event.preventDefault()}
          className={cn(
            "inline-flex",
            "h-8",
            "min-w-11",
            "flex-1",
            "items-center",
            "justify-center",
            "rounded-md",
            "px-3",
            "text-muted-foreground",
            "transition-colors",
            "hover:bg-muted",
            "hover:text-foreground",
            {
              "bg-muted text-foreground": active,
            }
          )}
        >
          {children}
        </button>
      </TooltipTrigger>
      <TooltipContent
        side="top"
        sideOffset={8}
        className={cn(
          "z-[60]",
          "border",
          "border-border",
          "bg-popover",
          "text-popover-foreground",
          "shadow-md",
          "[&_svg]:bg-popover",
          "[&_svg]:fill-popover"
        )}
      >
        {label}
      </TooltipContent>
    </Tooltip>
  );
}

Header.displayName = "Header";
MobileHeader.displayName = "MobileHeader";
DesktopHeader.displayName = "DesktopHeader";
UserAccountMenu.displayName = "UserAccountMenu";
