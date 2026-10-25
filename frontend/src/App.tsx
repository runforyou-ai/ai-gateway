import {
  Authenticated,
  CanAccess,
  Refine,
  type I18nProvider,
} from "@refinedev/core";
import { DevtoolsPanel, DevtoolsProvider } from "@refinedev/devtools";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import {
  Boxes,
  KeyRound,
  LayoutDashboard,
  Server,
  Settings,
  Users,
} from "lucide-react";
import { useTranslation } from "react-i18next";

import {
  BrowserRouter,
  Navigate,
  Outlet,
  Route,
  Routes,
} from "react-router";
import routerProvider, {
  CatchAllNavigate,
  DocumentTitleHandler,
  UnsavedChangesNotifier,
} from "@refinedev/react-router";
import { accessControlProvider } from "./providers/access-control";
import { dataProvider } from "./providers/data";
import { authProvider } from "./providers/auth";
import { useNotificationProvider } from "./components/refine-ui/notification/use-notification-provider";
import { Toaster } from "./components/refine-ui/notification/toaster";
import { ThemeProvider } from "./components/refine-ui/theme/theme-provider";
import { Layout } from "./components/refine-ui/layout/layout";
import { ErrorComponent } from "./components/refine-ui/layout/error-component";
import { DashboardPage } from "./pages/dashboard";
import { PlaceholderListPage } from "./pages/placeholder-list";
import { ProfilePage } from "./pages/profile";
import { SettingsPage } from "./pages/settings";
import { AcceptInvitePage } from "./pages/accept-invite";
import { ForgotPasswordPage } from "./pages/forgot-password";
import { ResetPasswordPage } from "./pages/reset-password";
import { SignInPage } from "./pages/sign-in";
import { SignUpPage } from "./pages/sign-up";
import { UserCreatePage } from "./pages/users/create";
import { UserEditPage } from "./pages/users/edit";
import { UserListPage } from "./pages/users/list";
import { WelcomePage } from "./pages/welcome";
import "./i18n";
import "./App.css";

const CONSOLE_PREFIX = "/console";

function App() {
  const { t, i18n } = useTranslation();

  const i18nProvider: I18nProvider = {
    translate: (key: string, options?: Record<string, unknown>) =>
      String(t(key, options)),
    changeLocale: (lang: string) => i18n.changeLanguage(lang),
    getLocale: () => i18n.resolvedLanguage ?? i18n.language,
  };

  return (
    <BrowserRouter>
      <RefineKbarProvider>
        <ThemeProvider>
          <DevtoolsProvider>
            <Refine
              dataProvider={dataProvider}
              authProvider={authProvider}
              accessControlProvider={accessControlProvider}
              i18nProvider={i18nProvider}
              notificationProvider={useNotificationProvider()}
              routerProvider={routerProvider}
              resources={[
                {
                  name: "dashboard",
                  list: CONSOLE_PREFIX,
                  meta: {
                    label: t("resources.dashboard"),
                    icon: <LayoutDashboard className="h-4 w-4" />,
                  },
                },
                {
                  name: "users",
                  list: `${CONSOLE_PREFIX}/users`,
                  create: `${CONSOLE_PREFIX}/users/create`,
                  edit: `${CONSOLE_PREFIX}/users/edit/:id`,
                  meta: {
                    label: t("resources.users"),
                    icon: <Users className="h-4 w-4" />,
                  },
                },
                {
                  name: "models",
                  list: `${CONSOLE_PREFIX}/models`,
                  meta: {
                    label: t("resources.models"),
                    icon: <Boxes className="h-4 w-4" />,
                  },
                },
                {
                  name: "providers",
                  list: `${CONSOLE_PREFIX}/providers`,
                  meta: {
                    label: t("resources.providers"),
                    icon: <Server className="h-4 w-4" />,
                  },
                },
                {
                  name: "api-keys",
                  list: `${CONSOLE_PREFIX}/api-keys`,
                  meta: {
                    label: t("resources.api-keys"),
                    icon: <KeyRound className="h-4 w-4" />,
                  },
                },
                {
                  name: "settings",
                  list: `${CONSOLE_PREFIX}/settings`,
                  meta: {
                    label: t("resources.settings"),
                    icon: <Settings className="h-4 w-4" />,
                  },
                },
                {
                  name: "profile",
                  list: `${CONSOLE_PREFIX}/profile`,
                  meta: {
                    label: t("profile.title"),
                    hide: true,
                  },
                },
              ]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                projectId: "V6ZxvK-SG6KKs-ZwwEIB",
                title: {
                  text: "AI Gateway",
                },
              }}
            >
              <Routes>
                <Route index element={<WelcomePage />} />
                <Route path="/login" element={<SignInPage />} />
                <Route path="/register" element={<SignUpPage />} />
                <Route path="/forgot-password" element={<ForgotPasswordPage />} />
                <Route path="/reset-password" element={<ResetPasswordPage />} />
                <Route path="/accept-invite" element={<AcceptInvitePage />} />
                <Route
                  path={CONSOLE_PREFIX}
                  element={
                    <Authenticated
                      key="authenticated"
                      fallback={<CatchAllNavigate to="/login" />}
                    >
                      <Layout>
                        <Outlet />
                      </Layout>
                    </Authenticated>
                  }
                >
                  <Route index element={<DashboardPage />} />
                  <Route
                    path="users"
                    element={
                      <CanAccess
                        resource="users"
                        action="list"
                        fallback={<Navigate to={CONSOLE_PREFIX} replace />}
                      >
                        <Outlet />
                      </CanAccess>
                    }
                  >
                    <Route index element={<UserListPage />} />
                    <Route path="create" element={<UserCreatePage />} />
                    <Route path="edit/:id" element={<UserEditPage />} />
                  </Route>
                  <Route path="models" element={<PlaceholderListPage />} />
                  <Route path="providers" element={<PlaceholderListPage />} />
                  <Route path="api-keys" element={<PlaceholderListPage />} />
                  <Route
                    path="settings"
                    element={
                      <CanAccess
                        resource="settings"
                        action="list"
                        fallback={<Navigate to={CONSOLE_PREFIX} replace />}
                      >
                        <SettingsPage />
                      </CanAccess>
                    }
                  />
                  <Route path="profile" element={<ProfilePage />} />
                  <Route path="*" element={<ErrorComponent />} />
                </Route>
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
              <Toaster />
              <RefineKbar />
              <UnsavedChangesNotifier />
              <DocumentTitleHandler />
            </Refine>
            <DevtoolsPanel />
          </DevtoolsProvider>
        </ThemeProvider>
      </RefineKbarProvider>
    </BrowserRouter>
  );
}

export default App;
