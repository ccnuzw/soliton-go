<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import {
  api,
  type FieldConfig,
  type FieldType,
  type GenerationResult,
  type DddListResponse,
} from '../api'
import { showSuccess, showError } from '../toast'

const loadingDomains = ref(false)
const domains = ref<string[]>([])
const selectedDomain = ref('')
const fieldTypes = ref<FieldType[]>([])
const dddList = ref<DddListResponse | null>(null)
const dddLoading = ref(false)
const selectedValueObjectName = ref('')
const selectedSpecName = ref('')
const selectedPolicyName = ref('')
const selectedEventName = ref('')
const selectedHandlerName = ref('')
const diffVisible = ref(false)
const diffTitle = ref('')
const diffExisting = ref('')
const diffPreview = ref('')
const diffFileName = ref('')
const diffHint = ref('')
const diffStatus = ref('')
const diffLoading = ref(false)
const valueObjectBatch = ref('')
const specBatch = ref('')
const policyBatch = ref('')
const eventBatch = ref('')

const valueObjectLoading = ref(false)
const valueObjectResult = ref<GenerationResult | null>(null)
const specLoading = ref(false)
const specResult = ref<GenerationResult | null>(null)
const policyLoading = ref(false)
const policyResult = ref<GenerationResult | null>(null)
const eventLoading = ref(false)
const eventResult = ref<GenerationResult | null>(null)

type ValueObjectState = {
  name: string
  fields: FieldConfig[]
  force: boolean
}

type SpecificationState = {
  name: string
  target: string
  force: boolean
}

type PolicyState = {
  name: string
  target: string
  force: boolean
}

type EventFlowState = {
  name: string
  topic: string
  fields: FieldConfig[]
  generateEvent: boolean
  generateHandler: boolean
  handlerTopic: string
  eventForce: boolean
  handlerForce: boolean
}

const valueObject = ref<ValueObjectState>({
  name: '',
  fields: [{ name: 'value', type: 'string', comment: '', enum_values: [] as string[] }],
  force: false,
})

const specification = ref<SpecificationState>({
  name: '',
  target: '',
  force: false,
})

const policy = ref<PolicyState>({
  name: '',
  target: '',
  force: false,
})

const eventFlow = ref<EventFlowState>({
  name: '',
  topic: '',
  fields: [{ name: 'user_id', type: 'uuid', comment: '', enum_values: [] as string[] }],
  generateEvent: true,
  generateHandler: true,
  handlerTopic: '',
  eventForce: false,
  handlerForce: false,
})

const domainHint = computed(() => {
  if (!selectedDomain.value.trim()) {
    return '请选择或输入已有领域名称'
  }
  const selectedLower = selectedDomain.value.trim().toLowerCase()
  const exists = domains.value.some((d) => d.toLowerCase() === selectedLower)
  if (domains.value.length > 0 && !exists) {
    return '未在项目中找到该领域，请先生成领域模块'
  }
  return '可选领域来自已生成模块'
})

onMounted(async () => {
  try {
    const res = await api.getFieldTypes()
    fieldTypes.value = res.types
    await loadDomains()
    await loadDddList()
  } catch (e) {
    console.error(e)
  }
})

watch(selectedDomain, async () => {
  selectedValueObjectName.value = ''
  selectedSpecName.value = ''
  selectedPolicyName.value = ''
  selectedEventName.value = ''
  selectedHandlerName.value = ''
  await loadDddList()
})

async function loadDomains() {
  loadingDomains.value = true
  try {
    const res = await api.listDomains()
    domains.value = res.domains.map((d) => d.name)
    if (!selectedDomain.value && domains.value.length > 0) {
      const first = domains.value[0]
      if (first) {
        selectedDomain.value = first
      }
    }
  } catch (e) {
    console.error('Failed to load domains:', e)
  } finally {
    loadingDomains.value = false
  }
}

async function loadDddList() {
  const domain = selectedDomain.value.trim()
  if (!domain) {
    dddList.value = null
    return
  }
  dddLoading.value = true
  try {
    dddList.value = await api.listDddComponents(domain)
  } catch (e) {
    console.error('Failed to load DDD list:', e)
  } finally {
    dddLoading.value = false
  }
}

function ensureDomain(): string | null {
  const domain = selectedDomain.value.trim()
  if (!domain) {
    showError('请先选择领域')
    return null
  }
  return domain
}

function addField(target: FieldConfig[]) {
  target.push({ name: '', type: 'string', comment: '', enum_values: [] })
}

function removeField(target: FieldConfig[], index: number) {
  if (target.length > 1) {
    target.splice(index, 1)
  }
}

function updateEnumValues(field: FieldConfig, value: string) {
  field.enum_values = value.split('|').map((v) => v.trim()).filter(Boolean)
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    new: '新建',
    overwrite: '覆盖',
    skip: '跳过',
    error: '错误',
  }
  return map[status] || status
}

function closeDiffModal() {
  diffVisible.value = false
}

function openDiffModal(
  title: string,
  existing: string,
  preview: string,
  fileName: string,
  status: string,
  hint = ''
) {
  diffTitle.value = title
  diffExisting.value = existing
  diffPreview.value = preview
  diffFileName.value = fileName
  diffStatus.value = status
  diffHint.value = hint
  diffVisible.value = true
}

function normalizeFieldEntry(entry: any): FieldConfig | null {
  if (!entry) return null
  if (typeof entry === 'string') {
    return parseFieldLine(entry)
  }
  if (typeof entry !== 'object') return null
  if (!entry.name || !entry.type) return null
  const enumValues =
    Array.isArray(entry.enum_values) ? entry.enum_values : Array.isArray(entry.enumValues) ? entry.enumValues : []
  return {
    name: String(entry.name),
    type: String(entry.type),
    comment: entry.comment ? String(entry.comment) : '',
    enum_values: enumValues,
  }
}

function parseFieldLine(line: string): FieldConfig | null {
  const trimmed = line.trim()
  if (!trimmed) return null

  let name = ''
  let fieldType = ''
  let comment = ''

  if (trimmed.includes(',')) {
    const parts = trimmed.split(',')
    name = (parts[0] || '').trim()
    fieldType = (parts[1] || '').trim()
    comment = parts.slice(2).join(',').trim()
  } else if (trimmed.includes('\t')) {
    const parts = trimmed.split('\t')
    name = (parts[0] || '').trim()
    fieldType = (parts[1] || '').trim()
    comment = parts.slice(2).join(' ').trim()
  } else {
    const match = trimmed.match(/^(\S+)\s+(\S+)\s*(.*)$/)
    if (!match) return null
    name = match[1] || ''
    fieldType = match[2] || ''
    comment = match[3] || ''
  }

  if (!name || !fieldType) return null
  let enumValues: string[] = []
  if (fieldType.startsWith('enum:')) {
    enumValues = fieldType
      .slice('enum:'.length)
      .split('|')
      .map((v) => v.trim())
      .filter(Boolean)
    fieldType = 'enum'
  }
  return {
    name,
    type: fieldType,
    comment: comment || '',
    enum_values: enumValues,
  }
}

function sanitizeFields(fields: FieldConfig[]): FieldConfig[] {
  return fields.map((field) => ({
    name: field.name || '',
    type: field.type || 'string',
    comment: field.comment || '',
    enum_values: field.enum_values ? [...field.enum_values] : [],
  }))
}

function normalizeFields(fields: FieldConfig[] | undefined, fallback: FieldConfig[]): FieldConfig[] {
  if (!fields || fields.length === 0) return fallback
  return sanitizeFields(fields)
}

function parseFieldBatchInput(input: string) {
  const trimmed = input.trim()
  if (!trimmed) return null
  try {
    const parsed = JSON.parse(trimmed)
    if (Array.isArray(parsed)) {
      const fields = parsed.map(normalizeFieldEntry).filter(Boolean) as FieldConfig[]
      return fields.length ? { fields } : null
    }
    if (parsed && typeof parsed === 'object') {
      const rawFields = Array.isArray(parsed.fields) ? parsed.fields : []
      const fields = rawFields.map(normalizeFieldEntry).filter(Boolean) as FieldConfig[]
      return {
        name: parsed.name || parsed.event_name,
        topic: parsed.topic,
        handlerTopic: parsed.handler_topic || parsed.handlerTopic,
        generateEvent: parsed.generateEvent,
        generateHandler: parsed.generateHandler,
        eventForce: parsed.eventForce,
        handlerForce: parsed.handlerForce,
        fields,
      }
    }
  } catch (e) {
    // Fall back to line parsing.
  }

  const fields = trimmed
    .split(/\r?\n/)
    .map(parseFieldLine)
    .filter(Boolean) as FieldConfig[]
  if (!fields.length) return null
  return { fields }
}

function parseNameTargetInput(input: string) {
  const trimmed = input.trim()
  if (!trimmed) return null
  try {
    const parsed = JSON.parse(trimmed)
    if (parsed && typeof parsed === 'object') {
      return {
        name: parsed.name,
        target: parsed.target,
        force: parsed.force,
      }
    }
  } catch (e) {
    // Fall back to line parsing.
  }

  const firstLine = trimmed.split(/\r?\n/).find((line) => line.trim())
  if (!firstLine) return null
  if (firstLine.includes(',')) {
    const parts = firstLine.split(',')
    return { name: (parts[0] || '').trim(), target: (parts[1] || '').trim() }
  }
  const match = firstLine.match(/^(\S+)\s+(\S+)\s*$/)
  if (match) {
    return { name: match[1] || '', target: match[2] || '' }
  }
  return { name: firstLine.trim(), target: '' }
}

function applyValueObjectBatch() {
  const parsed = parseFieldBatchInput(valueObjectBatch.value)
  if (!parsed) {
    showError('请输入有效的字段 JSON 或行列表')
    return
  }
  if (parsed.name) {
    valueObject.value.name = String(parsed.name)
  }
  if (parsed.fields?.length) {
    valueObject.value.fields = sanitizeFields(parsed.fields)
  }
  showSuccess('已导入字段')
}

function exportValueObjectBatch() {
  valueObjectBatch.value = JSON.stringify(
    {
      name: valueObject.value.name,
      fields: valueObject.value.fields,
    },
    null,
    2
  )
}

function applySpecBatch() {
  const parsed = parseNameTargetInput(specBatch.value)
  if (!parsed) {
    showError('请输入有效的 JSON 或 Name,Target')
    return
  }
  if (parsed.name) {
    specification.value.name = String(parsed.name)
  }
  if (parsed.target !== undefined) {
    specification.value.target = parsed.target || ''
  }
  if (typeof parsed.force === 'boolean') {
    specification.value.force = parsed.force
  }
  showSuccess('已导入 Specification 配置')
}

function exportSpecBatch() {
  specBatch.value = JSON.stringify(
    {
      name: specification.value.name,
      target: specification.value.target,
      force: specification.value.force,
    },
    null,
    2
  )
}

function applyPolicyBatch() {
  const parsed = parseNameTargetInput(policyBatch.value)
  if (!parsed) {
    showError('请输入有效的 JSON 或 Name,Target')
    return
  }
  if (parsed.name) {
    policy.value.name = String(parsed.name)
  }
  if (parsed.target !== undefined) {
    policy.value.target = parsed.target || ''
  }
  if (typeof parsed.force === 'boolean') {
    policy.value.force = parsed.force
  }
  showSuccess('已导入 Policy 配置')
}

function exportPolicyBatch() {
  policyBatch.value = JSON.stringify(
    {
      name: policy.value.name,
      target: policy.value.target,
      force: policy.value.force,
    },
    null,
    2
  )
}

function applyEventBatch() {
  const parsed = parseFieldBatchInput(eventBatch.value)
  if (!parsed) {
    showError('请输入有效的字段 JSON 或行列表')
    return
  }
  if (parsed.name) {
    eventFlow.value.name = String(parsed.name)
  }
  if (parsed.topic !== undefined) {
    eventFlow.value.topic = parsed.topic || ''
  }
  if (parsed.handlerTopic !== undefined) {
    eventFlow.value.handlerTopic = parsed.handlerTopic || ''
  }
  if (typeof parsed.generateEvent === 'boolean') {
    eventFlow.value.generateEvent = parsed.generateEvent
  }
  if (typeof parsed.generateHandler === 'boolean') {
    eventFlow.value.generateHandler = parsed.generateHandler
  }
  if (typeof parsed.eventForce === 'boolean') {
    eventFlow.value.eventForce = parsed.eventForce
  }
  if (typeof parsed.handlerForce === 'boolean') {
    eventFlow.value.handlerForce = parsed.handlerForce
  }
  if (parsed.fields?.length) {
    eventFlow.value.fields = sanitizeFields(parsed.fields)
    eventFlow.value.generateEvent = true
  }
  showSuccess('已导入 Event 配置')
}

function exportEventBatch() {
  eventBatch.value = JSON.stringify(
    {
      name: eventFlow.value.name,
      topic: eventFlow.value.topic,
      handler_topic: eventFlow.value.handlerTopic,
      generateEvent: eventFlow.value.generateEvent,
      generateHandler: eventFlow.value.generateHandler,
      eventForce: eventFlow.value.eventForce,
      handlerForce: eventFlow.value.handlerForce,
      fields: eventFlow.value.fields,
    },
    null,
    2
  )
}

function applyRenameResult(itemType: string, oldName: string, newName: string) {
  if (itemType === 'valueobject') {
    if (selectedValueObjectName.value === oldName) selectedValueObjectName.value = newName
    if (valueObject.value.name === oldName) valueObject.value.name = newName
  }
  if (itemType === 'spec') {
    if (selectedSpecName.value === oldName) selectedSpecName.value = newName
    if (specification.value.name === oldName) specification.value.name = newName
  }
  if (itemType === 'policy') {
    if (selectedPolicyName.value === oldName) selectedPolicyName.value = newName
    if (policy.value.name === oldName) policy.value.name = newName
  }
  if (itemType === 'event') {
    if (selectedEventName.value === oldName) selectedEventName.value = newName
    if (eventFlow.value.name === oldName) eventFlow.value.name = newName
  }
  if (itemType === 'event_handler') {
    if (selectedHandlerName.value === oldName) selectedHandlerName.value = newName
    if (eventFlow.value.name === oldName) eventFlow.value.name = newName
  }
}

function applyDeleteResult(itemType: string, name: string) {
  if (itemType === 'valueobject' && selectedValueObjectName.value === name) {
    selectedValueObjectName.value = ''
  }
  if (itemType === 'spec' && selectedSpecName.value === name) {
    selectedSpecName.value = ''
  }
  if (itemType === 'policy' && selectedPolicyName.value === name) {
    selectedPolicyName.value = ''
  }
  if (itemType === 'event' && selectedEventName.value === name) {
    selectedEventName.value = ''
  }
  if (itemType === 'event_handler' && selectedHandlerName.value === name) {
    selectedHandlerName.value = ''
  }
}

async function renameDddItem(itemType: string, name: string, label: string) {
  const domain = ensureDomain()
  if (!domain) return
  if (!name) {
    showError(`请选择需要重命名的 ${label}`)
    return
  }
  const newName = window.prompt(`请输入新的 ${label} 名称`, name)
  if (!newName || !newName.trim()) {
    return
  }
  if (newName.trim() === name) {
    showError('新名称与原名称一致')
    return
  }

  dddLoading.value = true
  try {
    await api.renameDddItem({
      domain,
      type: itemType,
      name,
      new_name: newName.trim(),
      force: false,
    })
    applyRenameResult(itemType, name, newName.trim())
    await loadDddList()
    showSuccess(`${label} 已重命名`)
  } catch (e: any) {
    if (e?.message && e.message.includes('target file exists')) {
      const confirmOverwrite = window.confirm('目标文件已存在，是否覆盖？')
      if (!confirmOverwrite) {
        return
      }
      try {
        await api.renameDddItem({
          domain,
          type: itemType,
          name,
          new_name: newName.trim(),
          force: true,
        })
        applyRenameResult(itemType, name, newName.trim())
        await loadDddList()
        showSuccess(`${label} 已覆盖重命名`)
      } catch (err: any) {
        showError(err?.message || '重命名失败')
      }
    } else {
      showError(e?.message || '重命名失败')
    }
  } finally {
    dddLoading.value = false
  }
}

async function deleteDddItem(itemType: string, name: string, label: string) {
  const domain = ensureDomain()
  if (!domain) return
  if (!name) {
    showError(`请选择需要删除的 ${label}`)
    return
  }
  const confirmDelete = window.confirm(`确认删除 ${label} ${name} 吗？`)
  if (!confirmDelete) return

  dddLoading.value = true
  try {
    await api.deleteDddItem({
      domain,
      type: itemType,
      name,
    })
    applyDeleteResult(itemType, name)
    await loadDddList()
    showSuccess(`${label} 已删除`)
  } catch (e: any) {
    showError(e?.message || '删除失败')
  } finally {
    dddLoading.value = false
  }
}

async function fetchExistingSource(domain: string, itemType: string, name: string) {
  try {
    const res = await api.getDddSource(domain, itemType, name)
    return { content: res.content, file: res.file, found: true }
  } catch (e) {
    return { content: '', file: '', found: false }
  }
}

function pickPreviewFile(result: GenerationResult) {
  const file = result.files.find((f) => f.content) || result.files[0]
  if (!file || !file.content) {
    throw new Error('预览内容为空')
  }
  return {
    content: file.content,
    file: file.path.split('/').pop() || file.path,
    status: file.status,
  }
}

async function showValueObjectDiff() {
  const domain = ensureDomain()
  if (!domain) return
  if (!valueObject.value.name.trim()) {
    showError('请输入 Value Object 名称')
    return
  }
  diffLoading.value = true
  try {
    const fields = valueObject.value.fields.filter((f) => f.name.trim())
    const preview = await api.previewValueObject({
      domain,
      name: valueObject.value.name,
      fields,
      force: valueObject.value.force,
    })
    const previewFile = pickPreviewFile(preview)
    const existing = await fetchExistingSource(domain, 'valueobject', valueObject.value.name)
    openDiffModal(
      `Value Object · ${valueObject.value.name}`,
      existing.content,
      previewFile.content,
      previewFile.file,
      previewFile.status,
      existing.found ? '' : '当前文件不存在，将创建新文件'
    )
  } catch (e: any) {
    showError(e?.message || '对比失败')
  } finally {
    diffLoading.value = false
  }
}

async function showSpecDiff() {
  const domain = ensureDomain()
  if (!domain) return
  if (!specification.value.name.trim()) {
    showError('请输入 Specification 名称')
    return
  }
  diffLoading.value = true
  try {
    const preview = await api.previewSpecification({
      domain,
      name: specification.value.name,
      target: specification.value.target.trim(),
      force: specification.value.force,
    })
    const previewFile = pickPreviewFile(preview)
    const existing = await fetchExistingSource(domain, 'spec', specification.value.name)
    openDiffModal(
      `Specification · ${specification.value.name}`,
      existing.content,
      previewFile.content,
      previewFile.file,
      previewFile.status,
      existing.found ? '' : '当前文件不存在，将创建新文件'
    )
  } catch (e: any) {
    showError(e?.message || '对比失败')
  } finally {
    diffLoading.value = false
  }
}

async function showPolicyDiff() {
  const domain = ensureDomain()
  if (!domain) return
  if (!policy.value.name.trim()) {
    showError('请输入 Policy 名称')
    return
  }
  diffLoading.value = true
  try {
    const preview = await api.previewPolicy({
      domain,
      name: policy.value.name,
      target: policy.value.target.trim(),
      force: policy.value.force,
    })
    const previewFile = pickPreviewFile(preview)
    const existing = await fetchExistingSource(domain, 'policy', policy.value.name)
    openDiffModal(
      `Policy · ${policy.value.name}`,
      existing.content,
      previewFile.content,
      previewFile.file,
      previewFile.status,
      existing.found ? '' : '当前文件不存在，将创建新文件'
    )
  } catch (e: any) {
    showError(e?.message || '对比失败')
  } finally {
    diffLoading.value = false
  }
}

async function showEventDiff() {
  const domain = ensureDomain()
  if (!domain) return
  if (!eventFlow.value.name.trim()) {
    showError('请输入 Event 名称')
    return
  }
  if (!eventFlow.value.generateEvent) {
    showError('请勾选生成 Event 后再对比')
    return
  }
  diffLoading.value = true
  try {
    const fields = eventFlow.value.fields.filter((f) => f.name.trim())
    const preview = await api.previewEvent({
      domain,
      name: eventFlow.value.name,
      fields,
      topic: eventFlow.value.topic.trim(),
      force: eventFlow.value.eventForce,
    })
    const previewFile = pickPreviewFile(preview)
    const existing = await fetchExistingSource(domain, 'event', eventFlow.value.name)
    openDiffModal(
      `Event · ${eventFlow.value.name}`,
      existing.content,
      previewFile.content,
      previewFile.file,
      previewFile.status,
      existing.found ? '' : '当前文件不存在，将创建新文件'
    )
  } catch (e: any) {
    showError(e?.message || '对比失败')
  } finally {
    diffLoading.value = false
  }
}

async function showHandlerDiff() {
  const domain = ensureDomain()
  if (!domain) return
  if (!eventFlow.value.name.trim()) {
    showError('请输入 Event 名称')
    return
  }
  if (!eventFlow.value.generateHandler) {
    showError('请勾选生成 Handler 后再对比')
    return
  }
  diffLoading.value = true
  try {
    const preview = await api.previewEventHandler({
      domain,
      event_name: eventFlow.value.name,
      topic: (eventFlow.value.handlerTopic || eventFlow.value.topic).trim(),
      force: eventFlow.value.handlerForce,
    })
    const previewFile = pickPreviewFile(preview)
    const existing = await fetchExistingSource(domain, 'event_handler', eventFlow.value.name)
    openDiffModal(
      `Handler · ${eventFlow.value.name}`,
      existing.content,
      previewFile.content,
      previewFile.file,
      previewFile.status,
      existing.found ? '' : '当前文件不存在，将创建新文件'
    )
  } catch (e: any) {
    showError(e?.message || '对比失败')
  } finally {
    diffLoading.value = false
  }
}

function mergeResults(results: GenerationResult[]): GenerationResult {
  const merged: GenerationResult = { success: true, files: [] }
  const messages: string[] = []
  const errors: string[] = []

  results.forEach((res) => {
    merged.success = merged.success && res.success
    merged.files.push(...res.files)
    if (res.message) {
      messages.push(res.message)
    }
    if (res.errors && res.errors.length > 0) {
      errors.push(...res.errors)
    }
  })

  if (messages.length > 0) {
    merged.message = messages.join(' / ')
  }
  if (errors.length > 0) {
    merged.errors = errors
  }
  return merged
}

async function loadValueObjectDetail() {
  const domain = ensureDomain()
  if (!domain) return
  if (!selectedValueObjectName.value) {
    showError('请选择 Value Object')
    return
  }
  dddLoading.value = true
  try {
    const detail = await api.getDddDetail(domain, 'valueobject', selectedValueObjectName.value)
    valueObject.value.name = detail.name || selectedValueObjectName.value
    valueObject.value.fields = normalizeFields(detail.fields, [
      { name: 'value', type: 'string', comment: '', enum_values: [] },
    ])
    valueObject.value.force = true
    showSuccess('已加载 Value Object')
  } catch (e: any) {
    showError(e.message || '加载失败')
  } finally {
    dddLoading.value = false
  }
}

async function loadSpecDetail() {
  const domain = ensureDomain()
  if (!domain) return
  if (!selectedSpecName.value) {
    showError('请选择 Specification')
    return
  }
  dddLoading.value = true
  try {
    const detail = await api.getDddDetail(domain, 'spec', selectedSpecName.value)
    specification.value.name = detail.name || selectedSpecName.value
    specification.value.target = detail.target || ''
    specification.value.force = true
    showSuccess('已加载 Specification')
  } catch (e: any) {
    showError(e.message || '加载失败')
  } finally {
    dddLoading.value = false
  }
}

async function loadPolicyDetail() {
  const domain = ensureDomain()
  if (!domain) return
  if (!selectedPolicyName.value) {
    showError('请选择 Policy')
    return
  }
  dddLoading.value = true
  try {
    const detail = await api.getDddDetail(domain, 'policy', selectedPolicyName.value)
    policy.value.name = detail.name || selectedPolicyName.value
    policy.value.target = detail.target || ''
    policy.value.force = true
    showSuccess('已加载 Policy')
  } catch (e: any) {
    showError(e.message || '加载失败')
  } finally {
    dddLoading.value = false
  }
}

async function loadEventDetail() {
  const domain = ensureDomain()
  if (!domain) return
  if (!selectedEventName.value) {
    showError('请选择 Event')
    return
  }
  dddLoading.value = true
  try {
    const detail = await api.getDddDetail(domain, 'event', selectedEventName.value)
    eventFlow.value.name = detail.name || selectedEventName.value
    eventFlow.value.topic = detail.topic || ''
    eventFlow.value.fields = normalizeFields(detail.fields, [
      { name: 'user_id', type: 'uuid', comment: '', enum_values: [] },
    ])
    eventFlow.value.generateEvent = true
    eventFlow.value.generateHandler = false
    eventFlow.value.eventForce = true
    showSuccess('已加载 Event')
  } catch (e: any) {
    showError(e.message || '加载失败')
  } finally {
    dddLoading.value = false
  }
}

async function loadEventHandlerDetail() {
  const domain = ensureDomain()
  if (!domain) return
  if (!selectedHandlerName.value) {
    showError('请选择 Handler')
    return
  }
  dddLoading.value = true
  try {
    const detail = await api.getDddDetail(domain, 'event_handler', selectedHandlerName.value)
    eventFlow.value.name = detail.event_name || selectedHandlerName.value
    eventFlow.value.handlerTopic = detail.topic || ''
    eventFlow.value.generateEvent = false
    eventFlow.value.generateHandler = true
    eventFlow.value.handlerForce = true
    showSuccess('已加载 Handler')
  } catch (e: any) {
    showError(e.message || '加载失败')
  } finally {
    dddLoading.value = false
  }
}

async function previewValueObject() {
  const domain = ensureDomain()
  if (!domain) return
  if (!valueObject.value.name.trim()) {
    showError('请输入 Value Object 名称')
    return
  }

  valueObjectLoading.value = true
  try {
    const fields = valueObject.value.fields.filter((f) => f.name.trim())
    valueObjectResult.value = await api.previewValueObject({
      domain,
      name: valueObject.value.name,
      fields,
      force: valueObject.value.force,
    })
  } catch (e: any) {
    showError(e.message)
  } finally {
    valueObjectLoading.value = false
  }
}

async function generateValueObject() {
  const domain = ensureDomain()
  if (!domain) return
  if (!valueObject.value.name.trim()) {
    showError('请输入 Value Object 名称')
    return
  }

  valueObjectLoading.value = true
  try {
    const fields = valueObject.value.fields.filter((f) => f.name.trim())
    valueObjectResult.value = await api.generateValueObject({
      domain,
      name: valueObject.value.name,
      fields,
      force: valueObject.value.force,
    })
    if (valueObjectResult.value.success) {
      showSuccess(valueObjectResult.value.message || 'Value Object 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    valueObjectLoading.value = false
  }
}

async function previewSpec() {
  const domain = ensureDomain()
  if (!domain) return
  if (!specification.value.name.trim()) {
    showError('请输入 Specification 名称')
    return
  }

  specLoading.value = true
  try {
    specResult.value = await api.previewSpecification({
      domain,
      name: specification.value.name,
      target: specification.value.target.trim(),
      force: specification.value.force,
    })
  } catch (e: any) {
    showError(e.message)
  } finally {
    specLoading.value = false
  }
}

async function generateSpec() {
  const domain = ensureDomain()
  if (!domain) return
  if (!specification.value.name.trim()) {
    showError('请输入 Specification 名称')
    return
  }

  specLoading.value = true
  try {
    specResult.value = await api.generateSpecification({
      domain,
      name: specification.value.name,
      target: specification.value.target.trim(),
      force: specification.value.force,
    })
    if (specResult.value.success) {
      showSuccess(specResult.value.message || 'Specification 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    specLoading.value = false
  }
}

async function previewPolicy() {
  const domain = ensureDomain()
  if (!domain) return
  if (!policy.value.name.trim()) {
    showError('请输入 Policy 名称')
    return
  }

  policyLoading.value = true
  try {
    policyResult.value = await api.previewPolicy({
      domain,
      name: policy.value.name,
      target: policy.value.target.trim(),
      force: policy.value.force,
    })
  } catch (e: any) {
    showError(e.message)
  } finally {
    policyLoading.value = false
  }
}

async function generatePolicy() {
  const domain = ensureDomain()
  if (!domain) return
  if (!policy.value.name.trim()) {
    showError('请输入 Policy 名称')
    return
  }

  policyLoading.value = true
  try {
    policyResult.value = await api.generatePolicy({
      domain,
      name: policy.value.name,
      target: policy.value.target.trim(),
      force: policy.value.force,
    })
    if (policyResult.value.success) {
      showSuccess(policyResult.value.message || 'Policy 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    policyLoading.value = false
  }
}

async function previewEventFlow() {
  const domain = ensureDomain()
  if (!domain) return
  if (!eventFlow.value.name.trim()) {
    showError('请输入 Event 名称')
    return
  }
  if (!eventFlow.value.generateEvent && !eventFlow.value.generateHandler) {
    showError('请至少选择一个生成项')
    return
  }

  eventLoading.value = true
  try {
    const tasks: GenerationResult[] = []
    if (eventFlow.value.generateEvent) {
      const fields = eventFlow.value.fields.filter((f) => f.name.trim())
      const res = await api.previewEvent({
        domain,
        name: eventFlow.value.name,
        fields,
        topic: eventFlow.value.topic.trim(),
        force: eventFlow.value.eventForce,
      })
      tasks.push(res)
    }

    if (eventFlow.value.generateHandler) {
      const res = await api.previewEventHandler({
        domain,
        event_name: eventFlow.value.name,
        topic: (eventFlow.value.handlerTopic || eventFlow.value.topic).trim(),
        force: eventFlow.value.handlerForce,
      })
      tasks.push(res)
    }

    eventResult.value = mergeResults(tasks)
  } catch (e: any) {
    showError(e.message)
  } finally {
    eventLoading.value = false
  }
}

async function generateEventFlow() {
  const domain = ensureDomain()
  if (!domain) return
  if (!eventFlow.value.name.trim()) {
    showError('请输入 Event 名称')
    return
  }
  if (!eventFlow.value.generateEvent && !eventFlow.value.generateHandler) {
    showError('请至少选择一个生成项')
    return
  }

  eventLoading.value = true
  try {
    const tasks: GenerationResult[] = []
    if (eventFlow.value.generateEvent) {
      const fields = eventFlow.value.fields.filter((f) => f.name.trim())
      const res = await api.generateEvent({
        domain,
        name: eventFlow.value.name,
        fields,
        topic: eventFlow.value.topic.trim(),
        force: eventFlow.value.eventForce,
      })
      tasks.push(res)
    }

    if (eventFlow.value.generateHandler) {
      const res = await api.generateEventHandler({
        domain,
        event_name: eventFlow.value.name,
        topic: (eventFlow.value.handlerTopic || eventFlow.value.topic).trim(),
        force: eventFlow.value.handlerForce,
      })
      tasks.push(res)
    }

    eventResult.value = mergeResults(tasks)
    if (eventResult.value.success) {
      showSuccess(eventResult.value.message || 'Event 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    eventLoading.value = false
  }
}
</script>

<template>
  <div class="editor">
    <h1>🧩 领域增强 DDD</h1>
    <p class="subtitle">以中文为主，专业术语保留英文：Value Object / Specification / Policy / Event & Handler</p>

    <details class="help-tips">
      <summary>📖 使用说明 Usage Guide</summary>
      <div class="tips-content">
        <p><strong>适用范围：</strong>以下功能依赖已有领域模块，请先在「领域模块」中生成。</p>
        <p><strong>Value Object：</strong>领域值对象，用于表达不可变的业务概念。</p>
        <p><strong>Specification：</strong>领域规格，用于描述是否满足某个业务条件。</p>
        <p><strong>Policy：</strong>领域策略，用于约束或判断业务行为。</p>
        <p><strong>命名建议：</strong>使用 PascalCase，如 <code>EmailAddress</code>、<code>ActiveUserSpec</code></p>
        <p><strong>Event 与 Handler：</strong>可单独生成，也可在同一操作中组合生成。</p>
        <p><strong>Topic：</strong>留空时自动按领域生成，例如 <code>user.activated</code></p>
        <p><strong>回显与管理：</strong>可从已有列表加载已生成组件，支持重命名与删除。</p>
        <p><strong>Diff 对比：</strong>对比当前文件与预览结果，Handler 对比仅显示处理器文件。</p>
        <p><strong>批量导入：</strong>支持 JSON 或行格式（如 <code>name,type,comment</code>），枚举可用 <code>enum:a|b</code>。</p>
      </div>
    </details>

    <div class="domain-select">
      <label>
        选择领域 Domain *
        <span class="tooltip" data-tooltip="需要已存在的领域模块">ⓘ</span>
      </label>
      <div class="domain-input">
        <input v-model="selectedDomain" list="domain-list" placeholder="user / order / product" />
        <datalist id="domain-list">
          <option v-for="d in domains" :key="d" :value="d">{{ d }}</option>
        </datalist>
        <select v-model="selectedDomain" class="domain-select">
          <option value="">选择已有领域</option>
          <option v-for="d in domains" :key="d" :value="d">{{ d }}</option>
        </select>
        <button class="btn" @click="loadDomains" :disabled="loadingDomains">
          {{ loadingDomains ? '刷新中...' : '刷新领域' }}
        </button>
      </div>
      <span class="hint">{{ domainHint }}</span>
    </div>

    <div class="grid">
      <!-- Value Object -->
      <section class="card">
        <div class="card-header">
          <h2>值对象 Value Object</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">用于建模不可变概念（如金额、邮箱、地址）。</p>
        <div class="existing">
          <label>已有 Value Object</label>
          <div class="existing-row">
            <select v-model="selectedValueObjectName" class="existing-select">
              <option value="">选择...</option>
              <option v-for="item in (dddList?.value_objects || [])" :key="item.name" :value="item.name">
                {{ item.name }}
              </option>
            </select>
            <div class="existing-actions">
              <button class="btn" @click="loadValueObjectDetail" :disabled="!selectedValueObjectName || dddLoading">
                加载
              </button>
              <button
                class="btn"
                @click="renameDddItem('valueobject', selectedValueObjectName, 'Value Object')"
                :disabled="!selectedValueObjectName || dddLoading"
              >
                重命名
              </button>
              <button
                class="btn danger"
                @click="deleteDddItem('valueobject', selectedValueObjectName, 'Value Object')"
                :disabled="!selectedValueObjectName || dddLoading"
              >
                删除
              </button>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label>名称 Name *</label>
          <input v-model="valueObject.name" placeholder="EmailAddress" />
        </div>

        <div class="fields-section">
          <div class="section-header">
            <h3>字段 Fields</h3>
            <button class="btn-add" @click="addField(valueObject.fields)">+ 添加字段</button>
          </div>
          <div class="field-row" v-for="(field, index) in valueObject.fields" :key="index">
            <input v-model="field.name" placeholder="value" class="field-name" />
            <select v-model="field.type" class="field-type">
              <option v-for="t in fieldTypes" :key="t.type" :value="t.type">
                {{ t.type }} - {{ t.description }}
              </option>
            </select>
            <input v-if="field.type === 'enum'" :value="field.enum_values?.join('|')"
              @input="updateEnumValues(field, ($event.target as HTMLInputElement).value)"
              placeholder="basic|premium" class="field-enum" />
            <input v-model="field.comment" placeholder="字段备注" class="field-comment" />
            <button class="btn-remove" @click="removeField(valueObject.fields, index)"
              :disabled="valueObject.fields.length === 1">×</button>
          </div>
        </div>

        <div class="batch-import">
          <div class="batch-header">
            <div class="batch-title">批量导入/导出字段</div>
            <div class="batch-actions">
              <button class="btn" @click="applyValueObjectBatch">导入</button>
              <button class="btn" @click="exportValueObjectBatch">导出</button>
            </div>
          </div>
          <textarea
            v-model="valueObjectBatch"
            placeholder="支持 JSON 数组或行格式：name,type,comment"
          ></textarea>
          <span class="hint">示例：{"name":"EmailAddress","fields":[{"name":"value","type":"string","comment":"邮箱"}]}</span>
        </div>

        <label class="checkbox">
          <input type="checkbox" v-model="valueObject.force" />
          强制覆盖 Force
        </label>

        <div class="actions">
          <button class="btn ghost" @click="showValueObjectDiff" :disabled="valueObjectLoading || diffLoading">
            {{ diffLoading ? '对比中...' : '对比 Diff' }}
          </button>
          <button class="btn" @click="previewValueObject" :disabled="valueObjectLoading">
            {{ valueObjectLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generateValueObject" :disabled="valueObjectLoading">
            {{ valueObjectLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>

        <div class="result" v-if="valueObjectResult">
          <div class="result-header">
            <span>{{ valueObjectResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in valueObjectResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Specification -->
      <section class="card">
        <div class="card-header">
          <h2>规格 Specification</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">用于封装业务规则判断（可复用）。</p>
        <div class="existing">
          <label>已有 Specification</label>
          <div class="existing-row">
            <select v-model="selectedSpecName" class="existing-select">
              <option value="">选择...</option>
              <option v-for="item in (dddList?.specs || [])" :key="item.name" :value="item.name">
                {{ item.name }}
              </option>
            </select>
            <div class="existing-actions">
              <button class="btn" @click="loadSpecDetail" :disabled="!selectedSpecName || dddLoading">
                加载
              </button>
              <button
                class="btn"
                @click="renameDddItem('spec', selectedSpecName, 'Specification')"
                :disabled="!selectedSpecName || dddLoading"
              >
                重命名
              </button>
              <button
                class="btn danger"
                @click="deleteDddItem('spec', selectedSpecName, 'Specification')"
                :disabled="!selectedSpecName || dddLoading"
              >
                删除
              </button>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label>名称 Name *</label>
          <input v-model="specification.name" placeholder="ActiveUserSpec" />
        </div>
        <div class="form-group">
          <label>目标类型 Target（可选）</label>
          <input v-model="specification.target" :placeholder="selectedDomain || 'User'" />
          <span class="hint">留空表示 any</span>
        </div>
        <label class="checkbox">
          <input type="checkbox" v-model="specification.force" />
          强制覆盖 Force
        </label>
        <div class="batch-import">
          <div class="batch-header">
            <div class="batch-title">批量导入/导出配置</div>
            <div class="batch-actions">
              <button class="btn" @click="applySpecBatch">导入</button>
              <button class="btn" @click="exportSpecBatch">导出</button>
            </div>
          </div>
          <textarea v-model="specBatch" placeholder="JSON 或 Name,Target"></textarea>
        </div>
        <div class="actions">
          <button class="btn ghost" @click="showSpecDiff" :disabled="specLoading || diffLoading">
            {{ diffLoading ? '对比中...' : '对比 Diff' }}
          </button>
          <button class="btn" @click="previewSpec" :disabled="specLoading">
            {{ specLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generateSpec" :disabled="specLoading">
            {{ specLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>
        <div class="result" v-if="specResult">
          <div class="result-header">
            <span>{{ specResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in specResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Policy -->
      <section class="card">
        <div class="card-header">
          <h2>策略 Policy</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">用于表达业务策略或决策逻辑。</p>
        <div class="existing">
          <label>已有 Policy</label>
          <div class="existing-row">
            <select v-model="selectedPolicyName" class="existing-select">
              <option value="">选择...</option>
              <option v-for="item in (dddList?.policies || [])" :key="item.name" :value="item.name">
                {{ item.name }}
              </option>
            </select>
            <div class="existing-actions">
              <button class="btn" @click="loadPolicyDetail" :disabled="!selectedPolicyName || dddLoading">
                加载
              </button>
              <button
                class="btn"
                @click="renameDddItem('policy', selectedPolicyName, 'Policy')"
                :disabled="!selectedPolicyName || dddLoading"
              >
                重命名
              </button>
              <button
                class="btn danger"
                @click="deleteDddItem('policy', selectedPolicyName, 'Policy')"
                :disabled="!selectedPolicyName || dddLoading"
              >
                删除
              </button>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label>名称 Name *</label>
          <input v-model="policy.name" placeholder="PasswordPolicy" />
        </div>
        <div class="form-group">
          <label>目标类型 Target（可选）</label>
          <input v-model="policy.target" :placeholder="selectedDomain || 'User'" />
          <span class="hint">留空表示 any</span>
        </div>
        <label class="checkbox">
          <input type="checkbox" v-model="policy.force" />
          强制覆盖 Force
        </label>
        <div class="batch-import">
          <div class="batch-header">
            <div class="batch-title">批量导入/导出配置</div>
            <div class="batch-actions">
              <button class="btn" @click="applyPolicyBatch">导入</button>
              <button class="btn" @click="exportPolicyBatch">导出</button>
            </div>
          </div>
          <textarea v-model="policyBatch" placeholder="JSON 或 Name,Target"></textarea>
        </div>
        <div class="actions">
          <button class="btn ghost" @click="showPolicyDiff" :disabled="policyLoading || diffLoading">
            {{ diffLoading ? '对比中...' : '对比 Diff' }}
          </button>
          <button class="btn" @click="previewPolicy" :disabled="policyLoading">
            {{ policyLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generatePolicy" :disabled="policyLoading">
            {{ policyLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>
        <div class="result" v-if="policyResult">
          <div class="result-header">
            <span>{{ policyResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in policyResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Event + Handler -->
      <section class="card full">
        <div class="card-header">
          <h2>事件与处理器 Event & Handler</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">支持事件与处理器组合生成，自动注入模块与 EventBus。</p>
        <div class="existing">
          <label>已有 Event</label>
          <div class="existing-row">
            <select v-model="selectedEventName" class="existing-select">
              <option value="">选择...</option>
              <option v-for="item in (dddList?.events || [])" :key="item.name" :value="item.name">
                {{ item.name }}
              </option>
            </select>
            <div class="existing-actions">
              <button class="btn" @click="loadEventDetail" :disabled="!selectedEventName || dddLoading">
                加载
              </button>
              <button
                class="btn"
                @click="renameDddItem('event', selectedEventName, 'Event')"
                :disabled="!selectedEventName || dddLoading"
              >
                重命名
              </button>
              <button
                class="btn danger"
                @click="deleteDddItem('event', selectedEventName, 'Event')"
                :disabled="!selectedEventName || dddLoading"
              >
                删除
              </button>
            </div>
          </div>
        </div>
        <div class="existing">
          <label>已有 Handler</label>
          <div class="existing-row">
            <select v-model="selectedHandlerName" class="existing-select">
              <option value="">选择...</option>
              <option v-for="item in (dddList?.event_handlers || [])" :key="item.name" :value="item.name">
                {{ item.name }}
              </option>
            </select>
            <div class="existing-actions">
              <button class="btn" @click="loadEventHandlerDetail" :disabled="!selectedHandlerName || dddLoading">
                加载
              </button>
              <button
                class="btn"
                @click="renameDddItem('event_handler', selectedHandlerName, 'Handler')"
                :disabled="!selectedHandlerName || dddLoading"
              >
                重命名
              </button>
              <button
                class="btn danger"
                @click="deleteDddItem('event_handler', selectedHandlerName, 'Handler')"
                :disabled="!selectedHandlerName || dddLoading"
              >
                删除
              </button>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label>事件名称 Event *</label>
          <input v-model="eventFlow.name" placeholder="UserActivated" />
        </div>
        <div class="form-group">
          <label>事件 Topic（可选）</label>
          <input v-model="eventFlow.topic" placeholder="user.activated" />
          <span class="hint">留空自动生成 topic</span>
        </div>

        <div class="fields-section" v-if="eventFlow.generateEvent">
          <div class="section-header">
            <h3>事件字段 Event Fields</h3>
            <button class="btn-add" @click="addField(eventFlow.fields)">+ 添加字段</button>
          </div>
          <div class="field-row" v-for="(field, index) in eventFlow.fields" :key="index">
            <input v-model="field.name" placeholder="user_id" class="field-name" />
            <select v-model="field.type" class="field-type">
              <option v-for="t in fieldTypes" :key="t.type" :value="t.type">
                {{ t.type }} - {{ t.description }}
              </option>
            </select>
            <input v-if="field.type === 'enum'" :value="field.enum_values?.join('|')"
              @input="updateEnumValues(field, ($event.target as HTMLInputElement).value)"
              placeholder="created|deleted" class="field-enum" />
            <input v-model="field.comment" placeholder="字段备注" class="field-comment" />
            <button class="btn-remove" @click="removeField(eventFlow.fields, index)"
              :disabled="eventFlow.fields.length === 1">×</button>
          </div>
        </div>

        <div class="options-grid">
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.generateEvent" />
            生成 Event
          </label>
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.generateHandler" />
            生成 Handler
          </label>
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.eventForce" />
            Event 强制覆盖
          </label>
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.handlerForce" />
            Handler 强制覆盖
          </label>
        </div>

        <div class="form-group" v-if="eventFlow.generateHandler">
          <label>处理器 Topic（可选）</label>
          <input v-model="eventFlow.handlerTopic" placeholder="默认使用事件 Topic" />
          <span class="hint">留空时复用事件 Topic</span>
        </div>

        <div class="batch-import">
          <div class="batch-header">
            <div class="batch-title">批量导入/导出 Event 配置</div>
            <div class="batch-actions">
              <button class="btn" @click="applyEventBatch">导入</button>
              <button class="btn" @click="exportEventBatch">导出</button>
            </div>
          </div>
          <textarea v-model="eventBatch" placeholder="JSON 或行格式：name,type,comment"></textarea>
          <span class="hint">可携带 topic、handler_topic、generateEvent 等字段</span>
        </div>

        <div class="actions">
          <button class="btn ghost" @click="showEventDiff" :disabled="eventLoading || diffLoading">
            Event 对比
          </button>
          <button class="btn ghost" @click="showHandlerDiff" :disabled="eventLoading || diffLoading">
            Handler 对比
          </button>
          <button class="btn" @click="previewEventFlow" :disabled="eventLoading">
            {{ eventLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generateEventFlow" :disabled="eventLoading">
            {{ eventLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>

        <div class="result" v-if="eventResult">
          <div class="result-header">
            <span>{{ eventResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in eventResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
          <div class="message" v-if="eventResult.message">{{ eventResult.message }}</div>
        </div>
      </section>
    </div>

    <div v-if="diffVisible" class="diff-modal">
      <div class="diff-backdrop" @click="closeDiffModal"></div>
      <div class="diff-panel">
        <div class="diff-header">
          <div>
            <div class="diff-title">{{ diffTitle }}</div>
            <div class="diff-meta">
              <span v-if="diffFileName">文件：{{ diffFileName }}</span>
              <span v-if="diffStatus">状态：{{ getStatusText(diffStatus) }}</span>
              <span v-if="diffHint">{{ diffHint }}</span>
            </div>
          </div>
          <button class="btn" @click="closeDiffModal">关闭</button>
        </div>
        <div class="diff-body">
          <div class="diff-column">
            <h4>当前文件</h4>
            <pre>{{ diffExisting || '（当前文件不存在）' }}</pre>
          </div>
          <div class="diff-column">
            <h4>生成预览</h4>
            <pre>{{ diffPreview || '（暂无预览内容）' }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.editor {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px 40px;
}

.subtitle {
  color: var(--text-muted);
  margin-bottom: 16px;
}

.domain-select {
  margin: 16px 0 24px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  padding: 16px;
  border-radius: 12px;
}

.domain-input {
  display: flex;
  gap: 12px;
  align-items: center;
}

.domain-input input {
  flex: 1;
}

.domain-select {
  min-width: 180px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 18px;
}

.card.full {
  grid-column: 1 / -1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.tag {
  font-size: 0.75rem;
  padding: 4px 8px;
  background: rgba(99, 102, 241, 0.15);
  border-radius: 8px;
  color: var(--text);
}

.card-desc {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin-bottom: 12px;
}

.existing {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.existing-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.existing-actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.existing-select {
  flex: 1;
  min-width: 0;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.form-group label {
  font-weight: 500;
}

.form-group input,
.domain-input input,
select {
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 10px 12px;
  border-radius: 8px;
}

.hint {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.help-tips {
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid var(--primary);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.help-tips summary {
  cursor: pointer;
  font-weight: 500;
  color: var(--primary);
  user-select: none;
  list-style: none;
}

.help-tips summary::-webkit-details-marker {
  display: none;
}

.help-tips summary::before {
  content: '▶';
  display: inline-block;
  margin-right: 6px;
  font-size: 0.8em;
  transition: transform 0.2s;
}

.help-tips[open] summary::before {
  transform: rotate(90deg);
}

.tips-content {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.tips-content p {
  margin-bottom: 8px;
  color: var(--text-muted);
  line-height: 1.6;
}

.tips-content code {
  background: var(--bg-dark);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
}

.tooltip {
  display: inline-block;
  margin-left: 6px;
  color: var(--primary);
  cursor: help;
  font-size: 0.9em;
}

.tooltip:hover {
  color: var(--primary-dark);
}

.fields-section {
  margin-bottom: 12px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.field-row {
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr 1fr auto;
  gap: 8px;
  margin-bottom: 8px;
}

.field-name,
.field-type,
.field-enum,
.field-comment {
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 8px;
  border-radius: 8px;
}

.btn-add,
.btn,
.btn-remove,
.btn.primary {
  border: none;
  cursor: pointer;
  border-radius: 8px;
}

.btn {
  background: var(--bg-input);
  color: var(--text);
  padding: 8px 12px;
}

.btn.ghost {
  background: transparent;
  border: 1px solid var(--border);
}

.btn.danger {
  background: var(--error);
  color: white;
}

.btn.primary {
  background: var(--primary);
  color: white;
  padding: 8px 12px;
}

.btn-add {
  background: var(--primary);
  color: white;
  padding: 6px 10px;
}

.btn-remove {
  background: var(--error);
  color: white;
  padding: 0 10px;
}

.checkbox {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
  color: var(--text);
}

.actions {
  display: flex;
  gap: 10px;
  margin-top: 12px;
  flex-wrap: wrap;
}

.result {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 12px;
  margin-top: 16px;
}

.result-header {
  font-weight: 600;
  margin-bottom: 8px;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.file-status {
  padding: 2px 6px;
  border-radius: 6px;
  font-size: 0.75rem;
  text-transform: uppercase;
}

.file-status.new {
  background: var(--success);
  color: white;
}

.file-status.overwrite {
  background: var(--warning);
  color: white;
}

.file-status.skip {
  background: var(--bg-input);
  color: var(--text);
}

.file-status.error {
  background: var(--error);
  color: white;
}

.options-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-bottom: 12px;
}

.batch-import {
  background: rgba(99, 102, 241, 0.08);
  border: 1px dashed var(--border);
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 12px;
}

.batch-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  gap: 8px;
}

.batch-title {
  font-weight: 600;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.batch-import textarea {
  width: 100%;
  min-height: 120px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 10px 12px;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  resize: vertical;
}

.diff-modal {
  position: fixed;
  inset: 0;
  z-index: 50;
}

.diff-backdrop {
  position: absolute;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
}

.diff-panel {
  position: relative;
  z-index: 51;
  max-width: 1200px;
  margin: 60px auto;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 16px;
}

.diff-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
}

.diff-title {
  font-size: 1.1rem;
  font-weight: 600;
}

.diff-meta {
  display: flex;
  gap: 16px;
  font-size: 0.85rem;
  color: var(--text-muted);
  flex-wrap: wrap;
}

.diff-body {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.diff-column {
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 10px;
  min-height: 320px;
  display: flex;
  flex-direction: column;
}

.diff-column h4 {
  margin-bottom: 8px;
}

.diff-column pre {
  flex: 1;
  overflow: auto;
  white-space: pre-wrap;
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
}

@media (max-width: 960px) {
  .grid {
    grid-template-columns: 1fr;
  }

  .field-row {
    grid-template-columns: 1fr;
  }

  .domain-input {
    flex-direction: column;
    align-items: stretch;
  }

  .diff-body {
    grid-template-columns: 1fr;
  }

  .diff-panel {
    margin: 20px;
  }
}
</style>
