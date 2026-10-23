import { Authenticated, Refine } from "@refinedev/core";
import { DevtoolsPanel, DevtoolsProvider } from "@refinedev/devtools";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";

import {
  BrowserRouter,
  Navigate,
  Outlet,
  Route,
  Routes,
} from "react-router";
import routerProvider, {
  CatchAllNavigate,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router";
import { dataProvider } from "./providers/data";
import { authProvider } from "./providers/auth";
import { useNotificationProvider } from "./components/refine-ui/notification/use-notification-provider";
import { Toaster } from "./components/refine-ui/notification/toaster";
import { ThemeProvider } from "./components/refine-ui/theme/theme-provider";
import { EchoPage } from "./pages/echo";
import { LoginPage } from "./pages/login";
import { RegisterPage } from "./pages/register";
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
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                projectId: "V6ZxvK-SG6KKs-ZwwEIB",
              }}
            >
              <Routes>
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />
                <Route
                  element={
                    <Authenticated
                      key="authenticated"
                      fallback={<CatchAllNavigate to="/login" />}
                    >
                      <Outlet />
                    </Authenticated>
                  }
                >
                  <Route index element={<EchoPage />} />
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
