import { Authenticated, Refine } from "@refinedev/core";
import { DevtoolsPanel, DevtoolsProvider } from "@refinedev/devtools";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import {
  Boxes,
  KeyRound,
  LayoutDashboard,
  Radio,
  Server,
  Settings,
} from "lucide-react";

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
import { dataProvider } from "./providers/data";
import { authProvider } from "./providers/auth";
import { useNotificationProvider } from "./components/refine-ui/notification/use-notification-provider";
import { Toaster } from "./components/refine-ui/notification/toaster";
import { ThemeProvider } from "./components/refine-ui/theme/theme-provider";
import { Layout } from "./components/refine-ui/layout/layout";
import { ErrorComponent } from "./components/refine-ui/layout/error-component";
import { DashboardPage } from "./pages/dashboard";
import { EchoPage } from "./pages/echo";
import { PlaceholderListPage } from "./pages/placeholder-list";
import { SignInPage } from "./pages/sign-in";
import { SignUpPage } from "./pages/sign-up";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <RefineKbarProvider>
        <ThemeProvider>
          <DevtoolsProvider>
            <Refine
              dataProvider={dataProvider}
              authProvider={authProvider}
              notificationProvider={useNotificationProvider()}
              routerProvider={routerProvider}
              resources={[
                {
                  name: "dashboard",
                  list: "/",
                  meta: {
                    label: "概览",
                    icon: <LayoutDashboard className="h-4 w-4" />,
                  },
                },
                {
                  name: "models",
                  list: "/models",
                  meta: {
                    label: "模型",
                    icon: <Boxes className="h-4 w-4" />,
                  },
                },
                {
                  name: "providers",
                  list: "/providers",
                  meta: {
                    label: "上游凭据",
                    icon: <Server className="h-4 w-4" />,
                  },
                },
                {
                  name: "api-keys",
                  list: "/api-keys",
                  meta: {
                    label: "API Key",
                    icon: <KeyRound className="h-4 w-4" />,
                  },
                },
                {
                  name: "settings",
                  list: "/settings",
                  meta: {
                    label: "设置",
                    icon: <Settings className="h-4 w-4" />,
                  },
                },
                {
                  name: "echo",
                  list: "/echo",
                  meta: {
                    label: "Echo",
                    icon: <Radio className="h-4 w-4" />,
                  },
                },
              ]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                projectId: "V6ZxvK-SG6KKs-ZwwEIB",
                title: {
                  text: "Lite AI Gateway",
                },
              }}
            >
              <Routes>
                <Route path="/login" element={<SignInPage />} />
                <Route path="/register" element={<SignUpPage />} />
                <Route
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
                    path="/models"
                    element={
                      <PlaceholderListPage description="模型映射 CRUD 待接入。" />
                    }
                  />
                  <Route
                    path="/providers"
                    element={
                      <PlaceholderListPage description="上游供应商凭据 CRUD 待接入。" />
                    }
                  />
                  <Route
                    path="/api-keys"
                    element={
                      <PlaceholderListPage description="对外 API Key 管理待接入。" />
                    }
                  />
                  <Route
                    path="/settings"
                    element={
                      <PlaceholderListPage description="站点设置（如开放注册）待接入。" />
                    }
                  />
                  <Route path="/echo" element={<EchoPage />} />
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
