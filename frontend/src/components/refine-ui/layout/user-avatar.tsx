import { useGetIdentity } from "@refinedev/core";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

type User = {
  id: number;
  firstName: string;
  lastName: string;
  fullName: string;
  email: string;
  avatar?: string;
};

type UserAvatarProps = {
  className?: string;
  fallbackClassName?: string;
};

export function UserAvatar({ className, fallbackClassName }: UserAvatarProps) {
  const { data: user, isLoading: userIsLoading } = useGetIdentity<User>();

  if (userIsLoading || !user) {
    return (
      <Skeleton className={cn("size-8", "rounded-full", className)} />
    );
  }

  const { fullName, email, avatar } = user;

  return (
    <Avatar className={cn("size-8", "rounded-full", "select-none", className)}>
      {avatar && <AvatarImage src={avatar} alt={fullName || email} />}
      <AvatarFallback
        className={cn(
          "bg-black",
          "text-white",
          "text-sm",
          "font-medium",
          "select-none",
          fallbackClassName
        )}
      >
        {getInitials(fullName || email)}
      </AvatarFallback>
    </Avatar>
  );
}

const getInitials = (name = "") => {
  const names = name.split(" ");
  let initials = names[0].substring(0, 1).toUpperCase();

  if (names.length > 1) {
    initials += names[names.length - 1].substring(0, 1).toUpperCase();
  }
  return initials;
};

UserAvatar.displayName = "UserAvatar";
