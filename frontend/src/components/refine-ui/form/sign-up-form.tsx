"use client";

import { useState } from "react";

import {
  useRegister,
  useRefineOptions,
  useLink,
  useNotification,
  useTranslation,
} from "@refinedev/core";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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

type SignUpFormProps = {
  firstInstall?: boolean;
};

export const SignUpForm = ({ firstInstall = false }: SignUpFormProps) => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const { open } = useNotification();
  const Link = useLink();
  const { translate } = useTranslation();
  const { title } = useRefineOptions();
  const { mutate: register } = useRegister();

  const handleSignUp = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      open?.({
        type: "error",
        message: translate("pages.register.errors.passwordMismatch"),
        description: translate("pages.register.errors.passwordMismatchDesc"),
      });
      return;
    }

    register({
      name,
      email,
      password,
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
      <div className={cn("flex", "items-center", "justify-center", "gap-2")}>
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
              "text-green-600",
              "dark:text-green-400",
              "text-3xl",
              "font-semibold"
            )}
          >
            {translate("pages.register.title")}
          </CardTitle>
          <CardDescription
            className={cn("text-muted-foreground", "font-medium")}
          >
            {translate(
              firstInstall
                ? "pages.register.descriptionFirst"
                : "pages.register.description",
            )}
          </CardDescription>
        </CardHeader>

        <Separator />

        <CardContent className={cn("px-0")}>
          <form onSubmit={handleSignUp}>
            <div className={cn("flex", "flex-col", "gap-2")}>
              <Label htmlFor="name">
                {translate("pages.register.fields.name")}
              </Label>
              <Input
                id="name"
                type="text"
                autoComplete="nickname"
                required
                maxLength={64}
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>

            <div className={cn("flex", "flex-col", "gap-2", "mt-6")}>
              <Label htmlFor="email">
                {translate("pages.register.fields.email")}
              </Label>
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
              <Label htmlFor="password">
                {translate("pages.register.fields.password")}
              </Label>
              <InputPassword
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={8}
              />
            </div>

            <div
              className={cn("relative", "flex", "flex-col", "gap-2", "mt-6")}
            >
              <Label htmlFor="confirmPassword">
                {translate("pages.register.fields.confirmPassword")}
              </Label>
              <InputPassword
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                minLength={8}
              />
            </div>

            <Button
              type="submit"
              size="lg"
              className={cn(
                "w-full",
                "mt-6",
                "bg-green-600",
                "hover:bg-green-700",
                "text-white"
              )}
            >
              {translate("pages.register.buttons.submit")}
            </Button>
          </form>
        </CardContent>

        <Separator />

        <CardFooter>
          <div className={cn("w-full", "text-center text-sm")}>
            <span className={cn("text-sm", "text-muted-foreground")}>
              {translate("pages.register.buttons.haveAccount")}{" "}
            </span>
            <Link
              to="/login"
              className={cn(
                "text-blue-600",
                "dark:text-blue-400",
                "font-semibold",
                "underline"
              )}
            >
              {translate("pages.register.buttons.signin")}
            </Link>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
};

SignUpForm.displayName = "SignUpForm";
