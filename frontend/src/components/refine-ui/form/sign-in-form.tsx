"use client";

import { useState } from "react";

import { useLogin, useRefineOptions, useLink, useTranslation } from "@refinedev/core";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { InputPassword } from "@/components/refine-ui/form/input-password";
import { cn } from "@/lib/utils";

type SignInFormProps = {
  showSignUp?: boolean;
};

export const SignInForm = ({ showSignUp = false }: SignInFormProps) => {
  const [rememberMe, setRememberMe] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const Link = useLink();
  const { translate } = useTranslation();
  const { title } = useRefineOptions();
  const { mutate: login } = useLogin();

  const handleSignIn = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    login({
      email,
      password,
      remember: rememberMe,
    });
  };

  return (
    <div
      className={cn(
        "flex",
        "flex-col",
        "items-center",
        "justify-center",
        "px-6",
        "py-8",
        "min-h-svh"
      )}
    >
      <div className={cn("flex", "items-center", "justify-center")}>
        {title.icon && (
          <div
            className={cn("text-foreground", "[&>svg]:w-12", "[&>svg]:h-12")}
          >
            {title.icon}
          </div>
        )}
      </div>

      <Card className={cn("sm:w-[456px]", "p-12", "mt-6")}>
        <CardHeader className={cn("px-0")}>
          <CardTitle
            className={cn(
              "text-blue-600",
              "dark:text-blue-400",
              "text-3xl",
              "font-semibold"
            )}
          >
            {translate("pages.login.title")}
          </CardTitle>
          <CardDescription
            className={cn("text-muted-foreground", "font-medium")}
          >
            {translate("pages.login.description")}
          </CardDescription>
        </CardHeader>

        <Separator />

        <CardContent className={cn("px-0")}>
          <form onSubmit={handleSignIn}>
            <div className={cn("flex", "flex-col", "gap-2")}>
              <Label htmlFor="email">{translate("pages.login.fields.email")}</Label>
              <Input
                id="email"
                type="email"
                placeholder=""
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div
              className={cn("relative", "flex", "flex-col", "gap-2", "mt-6")}
            >
              <div className={cn("flex", "items-center", "justify-between")}>
                <Label htmlFor="password">
                  {translate("pages.login.fields.password")}
                </Label>
                <Link
                  to="/forgot-password"
                  className={cn(
                    "text-sm",
                    "text-muted-foreground",
                    "hover:text-foreground",
                    "transition-colors",
                  )}
                >
                  {translate("pages.login.buttons.forgotPassword")}
                </Link>
              </div>
              <InputPassword
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <div className={cn("flex items-center", "space-x-2", "mt-4")}>
              <Checkbox
                id="remember"
                checked={rememberMe}
                onCheckedChange={(checked) =>
                  setRememberMe(checked === "indeterminate" ? false : checked)
                }
              />
              <Label htmlFor="remember">
                {translate("pages.login.buttons.rememberMe")}
              </Label>
            </div>

            <Button type="submit" size="lg" className={cn("w-full", "mt-6")}>
              {translate("pages.login.buttons.submit")}
            </Button>
          </form>
        </CardContent>

        {showSignUp ? (
          <>
            <Separator />
            <CardFooter>
              <div className={cn("w-full", "text-center text-sm")}>
                <span className={cn("text-sm", "text-muted-foreground")}>
                  {translate("pages.login.buttons.noAccount")}{" "}
                </span>
                <Link
                  to="/register"
                  className={cn(
                    "text-green-600",
                    "dark:text-green-400",
                    "font-semibold",
                    "underline"
                  )}
                >
                  {translate("pages.login.buttons.signup")}
                </Link>
              </div>
            </CardFooter>
          </>
        ) : null}
      </Card>
    </div>
  );
};

SignInForm.displayName = "SignInForm";
