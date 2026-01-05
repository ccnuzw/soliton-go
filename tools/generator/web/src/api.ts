const API_BASE = '/api'

export interface FieldConfig {
  name: string
  type: string
  comment?: string
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

export interface ValueObjectConfig {
  domain: string
  name: string
  fields: FieldConfig[]
  force: boolean
}

export interface SpecificationConfig {
  domain: string
  name: string
  target?: string
  force: boolean
}

export interface PolicyConfig {
  domain: string
  name: string
  target?: string
  force: boolean
}

export interface EventConfig {
  domain: string
  name: string
  fields: FieldConfig[]
  topic?: string
  force: boolean
}

export interface EventHandlerConfig {
  domain: string
  event_name: string
  topic?: string
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

export interface MigrationLogEntry {
  time: string
  level: string
  step: string
  message: string
}

export interface MigrationResult {
  success: boolean
  message?: string
  logs: MigrationLogEntry[]
  duration_ms: number
  exit_code: number
  command: string
  started_at: string
  finished_at: string
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
  enum_values?: string[]
  comment?: string
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

export interface ServiceListItem {
  name: string
  methods: string[]
  type: 'domain_service' | 'cross_domain_service'
}

export interface ServiceMethodDetail {
  name: string
  camel_name: string
}

export interface ServiceDetail {
  name: string
  methods: ServiceMethodDetail[]
}

export interface ServiceDetectionResult {
  service_name: string
  domain_name: string
  domain_exists: boolean
  service_type: 'domain_service' | 'cross_domain_service'
  target_dir: string
  should_reuse_dto: boolean
  existing_dto_path?: string
  message: string
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

  deleteDomain: (name: string) =>
    request<{ success: boolean; message: string }>(`/domains/${name}`, {
      method: 'DELETE',
    }),

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

  listServices: () =>
    request<{ services: ServiceListItem[] }>('/services/list'),

  detectServiceType: (name: string) =>
    request<ServiceDetectionResult>(`/services/detect/${name}`),

  getServiceDetail: (name: string) =>
    request<ServiceDetail>(`/services/${name}`),

  deleteService: (name: string) =>
    request<{ success: boolean; message: string }>(`/services/${name}`, {
      method: 'DELETE',
    }),

  // Layout
  getLayout: () =>
    request<ProjectLayout>('/layout'),

  runMigration: (projectPath: string, autoTidy: boolean, timeoutSeconds: number) =>
    request<MigrationResult>('/projects/migrate', {
      method: 'POST',
      body: JSON.stringify({
        project_path: projectPath,
        auto_tidy: autoTidy,
        timeout_seconds: timeoutSeconds,
      }),
    }),

  // DDD
  generateValueObject: (config: ValueObjectConfig) =>
    request<GenerationResult>('/ddd/valueobjects', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewValueObject: (config: ValueObjectConfig) =>
    request<GenerationResult>('/ddd/valueobjects/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  generateSpecification: (config: SpecificationConfig) =>
    request<GenerationResult>('/ddd/specs', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewSpecification: (config: SpecificationConfig) =>
    request<GenerationResult>('/ddd/specs/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  generatePolicy: (config: PolicyConfig) =>
    request<GenerationResult>('/ddd/policies', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewPolicy: (config: PolicyConfig) =>
    request<GenerationResult>('/ddd/policies/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  generateEvent: (config: EventConfig) =>
    request<GenerationResult>('/ddd/events', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewEvent: (config: EventConfig) =>
    request<GenerationResult>('/ddd/events/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  generateEventHandler: (config: EventHandlerConfig) =>
    request<GenerationResult>('/ddd/event-handlers', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  previewEventHandler: (config: EventHandlerConfig) =>
    request<GenerationResult>('/ddd/event-handlers/preview', {
      method: 'POST',
      body: JSON.stringify(config),
    }),
}
