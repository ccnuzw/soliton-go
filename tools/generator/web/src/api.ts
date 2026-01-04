const API_BASE = '/api'

export interface FieldConfig {
  name: string
  type: string
  enum_values?: string[]
}

export interface ProjectConfig {
  name: string
  module_name?: string
  framework_version?: string
  framework_replace?: string
}

export interface DomainConfig {
  name: string
  fields: FieldConfig[]
  table_name?: string
  route_base?: string
  soft_delete: boolean
  wire: boolean
  force: boolean
}

export interface ServiceConfig {
  name: string
  methods: string[]
  force: boolean
}

export interface GeneratedFile {
  path: string
  status: string
  content?: string
}

export interface GenerationResult {
  success: boolean
  files: GeneratedFile[]
  errors?: string[]
  message?: string
}

export interface FieldType {
  type: string
  description: string
}

export interface ProjectLayout {
  found: boolean
  message?: string
  module_path?: string
  module_dir?: string
  internal_dir?: string
  domain_dir?: string
  app_dir?: string
  infra_dir?: string
  interfaces_dir?: string
}

export interface DomainListItem {
  name: string
  module_path: string
  fields: string[]
  has_files: boolean
}

export interface FieldDetail {
  name: string
  type: string
  is_enum: boolean
  gorm_tag: string
  json_tag: string
  snake_name: string
}

export interface DomainDetail {
  name: string
  fields: FieldDetail[]
  files: {
    entity: boolean
    repository: boolean
    events: boolean
  }
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  })
  
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Request failed')
  }
  
  return response.json()
}

export const api = {
  // Project
  initProject: (config: ProjectConfig) =>
    request<GenerationResult>('/projects/init', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewInitProject: (config: ProjectConfig) =>
    request<GenerationResult>('/projects/init/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  // Domain
  generateDomain: (config: DomainConfig) =>
    request<GenerationResult>('/domains', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewDomain: (config: DomainConfig) =>
    request<GenerationResult>('/domains/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  listDomains: () =>
    request<{ domains: DomainListItem[] }>('/domains/list'),

  getDomainDetail: (name: string) =>
    request<DomainDetail>(`/domains/${name}`),

  getFieldTypes: () =>
    request<{ types: FieldType[] }>('/field-types'),

  // Service
  generateService: (config: ServiceConfig) =>
    request<GenerationResult>('/services', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewService: (config: ServiceConfig) =>
    request<GenerationResult>('/services/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  // Layout
  getLayout: () =>
    request<ProjectLayout>('/layout'),
}
