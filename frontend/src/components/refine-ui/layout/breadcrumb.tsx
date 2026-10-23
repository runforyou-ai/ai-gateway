"use client";

import { Fragment, useMemo } from "react";
import { Home } from "lucide-react";
import {
  matchResourceFromRoute,
  useBreadcrumb,
  useLink,
  useResourceParams,
} from "@refinedev/core";
import {
  BreadcrumbSeparator as ShadcnBreadcrumbSeparator,
  BreadcrumbItem as ShadcnBreadcrumbItem,
  BreadcrumbList as ShadcnBreadcrumbList,
  BreadcrumbPage as ShadcnBreadcrumbPage,
  Breadcrumb as ShadcnBreadcrumb,
} from "@/components/ui/breadcrumb";

export function Breadcrumb() {
  const Link = useLink();
  const { breadcrumbs } = useBreadcrumb();
  const { resources } = useResourceParams();
  // 管理台挂在 /console，不能用模板默认的 "/"（会落到公开欢迎页）
  const rootRouteResource = matchResourceFromRoute("/console", resources);
  const homePath = rootRouteResource.matchedRoute ?? "/console";

  const breadCrumbItems = useMemo(() => {
    const list: {
      key: string;
      href: string;
      Component: React.ReactNode;
    }[] = [];

    list.push({
      key: "breadcrumb-item-home",
      href: homePath,
      Component: (
        <Link to={homePath}>
          <Home className="h-4 w-4" />
        </Link>
      ),
    });

    for (const { label, href } of breadcrumbs) {
      list.push({
        key: `breadcrumb-item-${label}`,
        href: href ?? "",
        Component: href ? <Link to={href}>{label}</Link> : <span>{label}</span>,
      });
    }

    return list;
  }, [breadcrumbs, Link, homePath]);

  return (
    <ShadcnBreadcrumb>
      <ShadcnBreadcrumbList>
        {breadCrumbItems.map((item, index) => {
          if (index === breadCrumbItems.length - 1) {
            return (
              <ShadcnBreadcrumbPage key={item.key}>
                {item.Component}
              </ShadcnBreadcrumbPage>
            );
          }

          return (
            <Fragment key={item.key}>
              <ShadcnBreadcrumbItem key={item.key}>
                {item.Component}
              </ShadcnBreadcrumbItem>
              <ShadcnBreadcrumbSeparator />
            </Fragment>
          );
        })}
      </ShadcnBreadcrumbList>
    </ShadcnBreadcrumb>
  );
}

Breadcrumb.displayName = "Breadcrumb";
