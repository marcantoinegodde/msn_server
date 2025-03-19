/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file was automatically generated by TanStack Router.
// You should NOT make any changes in this file as it will be overwritten.
// Additionally, you should also exclude this file from your linter and/or formatter to prevent it from being checked or modified.

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as LoginImport } from './routes/login'
import { Route as AuthRouteImport } from './routes/_auth/route'
import { Route as AuthLayoutRouteImport } from './routes/_auth/_layout/route'
import { Route as AuthLayoutIndexImport } from './routes/_auth/_layout/index'
import { Route as AuthLayoutStatusImport } from './routes/_auth/_layout/status'
import { Route as AuthLayoutDetailsImport } from './routes/_auth/_layout/details'

// Create/Update Routes

const LoginRoute = LoginImport.update({
  id: '/login',
  path: '/login',
  getParentRoute: () => rootRoute,
} as any)

const AuthRouteRoute = AuthRouteImport.update({
  id: '/_auth',
  getParentRoute: () => rootRoute,
} as any)

const AuthLayoutRouteRoute = AuthLayoutRouteImport.update({
  id: '/_layout',
  getParentRoute: () => AuthRouteRoute,
} as any)

const AuthLayoutIndexRoute = AuthLayoutIndexImport.update({
  id: '/',
  path: '/',
  getParentRoute: () => AuthLayoutRouteRoute,
} as any)

const AuthLayoutStatusRoute = AuthLayoutStatusImport.update({
  id: '/status',
  path: '/status',
  getParentRoute: () => AuthLayoutRouteRoute,
} as any)

const AuthLayoutDetailsRoute = AuthLayoutDetailsImport.update({
  id: '/details',
  path: '/details',
  getParentRoute: () => AuthLayoutRouteRoute,
} as any)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/_auth': {
      id: '/_auth'
      path: ''
      fullPath: ''
      preLoaderRoute: typeof AuthRouteImport
      parentRoute: typeof rootRoute
    }
    '/login': {
      id: '/login'
      path: '/login'
      fullPath: '/login'
      preLoaderRoute: typeof LoginImport
      parentRoute: typeof rootRoute
    }
    '/_auth/_layout': {
      id: '/_auth/_layout'
      path: ''
      fullPath: ''
      preLoaderRoute: typeof AuthLayoutRouteImport
      parentRoute: typeof AuthRouteImport
    }
    '/_auth/_layout/details': {
      id: '/_auth/_layout/details'
      path: '/details'
      fullPath: '/details'
      preLoaderRoute: typeof AuthLayoutDetailsImport
      parentRoute: typeof AuthLayoutRouteImport
    }
    '/_auth/_layout/status': {
      id: '/_auth/_layout/status'
      path: '/status'
      fullPath: '/status'
      preLoaderRoute: typeof AuthLayoutStatusImport
      parentRoute: typeof AuthLayoutRouteImport
    }
    '/_auth/_layout/': {
      id: '/_auth/_layout/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof AuthLayoutIndexImport
      parentRoute: typeof AuthLayoutRouteImport
    }
  }
}

// Create and export the route tree

interface AuthLayoutRouteRouteChildren {
  AuthLayoutDetailsRoute: typeof AuthLayoutDetailsRoute
  AuthLayoutStatusRoute: typeof AuthLayoutStatusRoute
  AuthLayoutIndexRoute: typeof AuthLayoutIndexRoute
}

const AuthLayoutRouteRouteChildren: AuthLayoutRouteRouteChildren = {
  AuthLayoutDetailsRoute: AuthLayoutDetailsRoute,
  AuthLayoutStatusRoute: AuthLayoutStatusRoute,
  AuthLayoutIndexRoute: AuthLayoutIndexRoute,
}

const AuthLayoutRouteRouteWithChildren = AuthLayoutRouteRoute._addFileChildren(
  AuthLayoutRouteRouteChildren,
)

interface AuthRouteRouteChildren {
  AuthLayoutRouteRoute: typeof AuthLayoutRouteRouteWithChildren
}

const AuthRouteRouteChildren: AuthRouteRouteChildren = {
  AuthLayoutRouteRoute: AuthLayoutRouteRouteWithChildren,
}

const AuthRouteRouteWithChildren = AuthRouteRoute._addFileChildren(
  AuthRouteRouteChildren,
)

export interface FileRoutesByFullPath {
  '': typeof AuthLayoutRouteRouteWithChildren
  '/login': typeof LoginRoute
  '/details': typeof AuthLayoutDetailsRoute
  '/status': typeof AuthLayoutStatusRoute
  '/': typeof AuthLayoutIndexRoute
}

export interface FileRoutesByTo {
  '': typeof AuthRouteRouteWithChildren
  '/login': typeof LoginRoute
  '/details': typeof AuthLayoutDetailsRoute
  '/status': typeof AuthLayoutStatusRoute
  '/': typeof AuthLayoutIndexRoute
}

export interface FileRoutesById {
  __root__: typeof rootRoute
  '/_auth': typeof AuthRouteRouteWithChildren
  '/login': typeof LoginRoute
  '/_auth/_layout': typeof AuthLayoutRouteRouteWithChildren
  '/_auth/_layout/details': typeof AuthLayoutDetailsRoute
  '/_auth/_layout/status': typeof AuthLayoutStatusRoute
  '/_auth/_layout/': typeof AuthLayoutIndexRoute
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths: '' | '/login' | '/details' | '/status' | '/'
  fileRoutesByTo: FileRoutesByTo
  to: '' | '/login' | '/details' | '/status' | '/'
  id:
    | '__root__'
    | '/_auth'
    | '/login'
    | '/_auth/_layout'
    | '/_auth/_layout/details'
    | '/_auth/_layout/status'
    | '/_auth/_layout/'
  fileRoutesById: FileRoutesById
}

export interface RootRouteChildren {
  AuthRouteRoute: typeof AuthRouteRouteWithChildren
  LoginRoute: typeof LoginRoute
}

const rootRouteChildren: RootRouteChildren = {
  AuthRouteRoute: AuthRouteRouteWithChildren,
  LoginRoute: LoginRoute,
}

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>()

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/_auth",
        "/login"
      ]
    },
    "/_auth": {
      "filePath": "_auth/route.tsx",
      "children": [
        "/_auth/_layout"
      ]
    },
    "/login": {
      "filePath": "login.tsx"
    },
    "/_auth/_layout": {
      "filePath": "_auth/_layout/route.tsx",
      "parent": "/_auth",
      "children": [
        "/_auth/_layout/details",
        "/_auth/_layout/status",
        "/_auth/_layout/"
      ]
    },
    "/_auth/_layout/details": {
      "filePath": "_auth/_layout/details.tsx",
      "parent": "/_auth/_layout"
    },
    "/_auth/_layout/status": {
      "filePath": "_auth/_layout/status.tsx",
      "parent": "/_auth/_layout"
    },
    "/_auth/_layout/": {
      "filePath": "_auth/_layout/index.tsx",
      "parent": "/_auth/_layout"
    }
  }
}
ROUTE_MANIFEST_END */
