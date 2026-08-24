package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SkillCatalogItem defines a standard skill definition for catalog generation
type SkillCatalogItem struct {
	ID           string
	Name         string
	Category     string
	Description  string
	ToolName     string
	ToolDesc     string
	Triggers     []string
	Keywords     []string
	Parameters   ToolParameters
	Instructions string
}

// GetUniversalSkillCatalog returns evaluation fixture schemas for router benchmarks.
// These are not live vendor integrations.
func GetUniversalSkillCatalog() []SkillCatalogItem {
	return []SkillCatalogItem{
		// === 1. Eval fixtures: documents and coding (not vendor-official) ===
		{
			ID:          "docx_document_architect",
			Name:        "Word Document Architect & Formatter",
			Category:    "Document Processing",
			Description: "Inspects, generates, formats and modifies Microsoft Word .docx documents, paragraph styles, headings and tables",
			ToolName:    "docx_manipulator",
			ToolDesc:    "Create and format Microsoft Word .docx documents",
			Triggers:    []string{"docx", "word", "document", "paragraph", "heading", "formatting", "word document"},
			Keywords:    []string{"docx", "document", "office", "style", "table", "msword"},
			Instructions: `# Word Document Architect & Formatter
Create and format Microsoft Word .docx documents with precise styling, headings, and table structures.`,
		},
		{
			ID:          "pdf_document_extractor",
			Name:        "PDF Text & Layout Extractor",
			Category:    "Document Processing",
			Description: "Extracts text, structured tables, metadata, and form fields from PDF documents with layout preservation",
			ToolName:    "pdf_extractor",
			ToolDesc:    "Extract text, tables, and structured data from PDF files",
			Triggers:    []string{"pdf", "ocr", "extract pdf", "acrobat", "pdf table", "pdf layout"},
			Keywords:    []string{"pdf", "extraction", "tables", "ocr", "document", "reader"},
			Instructions: `# PDF Text & Layout Extractor
Extract clean text, tables, and metadata from PDF files.`,
		},
		{
			ID:          "pptx_presentation_creator",
			Name:        "PowerPoint Deck Builder & Stylist",
			Category:    "Document Processing",
			Description: "Creates professional Microsoft PowerPoint .pptx slide presentations with layouts, bullets, and speaker notes",
			ToolName:    "pptx_deck_builder",
			ToolDesc:    "Generate and format PowerPoint .pptx presentation decks",
			Triggers:    []string{"pptx", "powerpoint", "slides", "deck", "presentation", "slide layout"},
			Keywords:    []string{"pptx", "presentation", "slideshow", "bullet", "keynote"},
			Instructions: `# PowerPoint Deck Builder & Stylist
Build structured slide decks with corporate templates, typography, and presenter notes.`,
		},
		{
			ID:          "xlsx_spreadsheet_modeler",
			Name:        "Excel Spreadsheet & Formula Modeler",
			Category:    "Document Processing",
			Description: "Builds, audits, and recalculates Excel .xlsx financial spreadsheets with SUMIFS, XLOOKUP, pivot tables, and styling",
			ToolName:    "xlsx_sheet_modeler",
			ToolDesc:    "Generate and audit Excel .xlsx financial spreadsheets and formulas",
			Triggers:    []string{"xlsx", "excel", "spreadsheet", "formula", "xlookup", "pivot table", "workbook"},
			Keywords:    []string{"xlsx", "excel", "sheet", "sumifs", "vlookup", "cells", "financial model"},
			Instructions: `# Excel Spreadsheet & Formula Modeler
Create robust spreadsheets with dynamic formulas, validation rules, and formatting.`,
		},
		{
			ID:          "theme_factory_stylist",
			Name:        "Design System & Theme Factory",
			Category:    "Design & UI",
			Description: "Generates cohesive UI color palettes, design tokens, typography pairings, and WCAG AA contrast matrices",
			ToolName:    "theme_factory_generate",
			ToolDesc:    "Generate design system color tokens and typography palettes",
			Triggers:    []string{"theme", "palette", "design token", "typography", "wcag", "color scheme", "ui design"},
			Keywords:    []string{"theme", "colors", "tokens", "styling", "css", "dark mode", "tailwind"},
			Instructions: `# Design System & Theme Factory
Generate modern, accessible UI color palettes and design tokens.`,
		},
		{
			ID:          "web_artifacts_builder",
			Name:        "Web Artifact & Prototype Builder",
			Category:    "Design & UI",
			Description: "Renders single-file interactive HTML, CSS, JavaScript, and SVG web prototypes and interactive dashboard widgets",
			ToolName:    "web_artifact_render",
			ToolDesc:    "Build and render interactive single-file web artifacts",
			Triggers:    []string{"web artifact", "prototype", "single file html", "interactive widget", "canvas app"},
			Keywords:    []string{"html", "css", "javascript", "artifact", "canvas", "frontend", "spa"},
			Instructions: `# Web Artifact & Prototype Builder
Build self-contained, responsive single-file web applications and interactive components.`,
		},
		{
			ID:          "mcp_server_builder",
			Name:        "MCP Server & Tool Manifest Builder",
			Category:    "AI & MCP",
			Description: "Scaffolds Model Context Protocol (MCP) servers, tool declarations, JSON-RPC 2.0 schemas, and STDIO/SSE handlers",
			ToolName:    "mcp_scaffold_server",
			ToolDesc:    "Scaffold and validate MCP Model Context Protocol servers and schemas",
			Triggers:    []string{"mcp", "model context protocol", "mcp server", "json-rpc", "tool manifest", "mcp tool"},
			Keywords:    []string{"mcp", "protocol", "context", "server", "tools", "manifest"},
			Instructions: `# MCP Server & Tool Manifest Builder
Scaffold production Model Context Protocol servers with schema validation.`,
		},
		{
			ID:          "prompt_engineer_distiller",
			Name:        "Prompt Optimization & Distiller",
			Category:    "AI & Prompting",
			Description: "Refines system instructions, few-shot prompt exemplars, DSPy-style signatures, and chain-of-thought constraints",
			ToolName:    "prompt_distill_optimize",
			ToolDesc:    "Optimize, distill, and benchmark LLM system prompts and few-shot examples",
			Triggers:    []string{"prompt engineering", "prompt optimization", "few-shot", "dspy", "system prompt", "distill prompt"},
			Keywords:    []string{"prompt", "instruction", "fewshot", "cot", "llm", "tuning"},
			Instructions: `# Prompt Optimization & Distiller
Refine prompts for maximum instruction adherence and minimum hallucination.`,
		},
		{
			ID:          "code_review_security_auditor",
			Name:        "Code Review & Security Hardener",
			Category:    "DevOps & Security",
			Description: "Audits source code for OWASP Top 10 vulnerabilities, SQL injection, XSS, buffer overflows, and architectural smells",
			ToolName:    "code_security_audit",
			ToolDesc:    "Audit source code for security vulnerabilities and architectural defects",
			Triggers:    []string{"code review", "security audit", "owasp", "sql injection", "vulnerability scan", "security hardener"},
			Keywords:    []string{"security", "audit", "cve", "owasp", "taint", "code review"},
			Instructions: `# Code Review & Security Hardener
Perform static code analysis to eliminate security vulnerabilities and performance bottlenecks.`,
		},
		{
			ID:          "sql_optimizer_query_tuner",
			Name:        "SQL Query Optimizer & Plan Tuner",
			Category:    "Database",
			Description: "Analyzes EXPLAIN ANALYZE execution plans, suggests index tuning, cost-based join reordering, and subquery flattening",
			ToolName:    "sql_plan_optimize",
			ToolDesc:    "Optimize SQL execution plans and index strategies",
			Triggers:    []string{"sql optimizer", "explain analyze", "query plan", "index tuning", "slow query", "query performance"},
			Keywords:    []string{"sql", "query", "optimizer", "explain", "index", "execution plan"},
			Instructions: `# SQL Query Optimizer & Plan Tuner
Analyze database execution trees and recommend targeted indexing and SQL rewrites.`,
		},

		// === 2. Google Gemini Official & Data Agent Kit Skills ===
		{
			ID:          "gemini_api_sdk_dev",
			Name:        "Gemini API & SDK Developer",
			Category:    "AI & Google Cloud",
			Description: "Develops applications using Google GenAI SDKs for Python, TypeScript, Go, function calling, and structured JSON outputs",
			ToolName:    "gemini_sdk_invoke",
			ToolDesc:    "Generate code and execute Google Gemini GenAI SDK applications",
			Triggers:    []string{"gemini sdk", "google genai", "structured output", "function calling", "gemini-flash", "gemini-pro"},
			Keywords:    []string{"gemini", "google", "genai", "sdk", "structured", "functioncall"},
			Instructions: `# Gemini API & SDK Developer
Build production GenAI applications using Google's official SDKs and structured schemas.`,
		},
		{
			ID:          "gemini_interactions_agent",
			Name:        "Gemini Interactions API Agent",
			Category:    "AI & Google Cloud",
			Description: "Orchestrates multi-turn conversational agents, background research loops, and multimodal interaction sessions",
			ToolName:    "gemini_interactions_call",
			ToolDesc:    "Execute multi-turn conversational turns with Gemini Interactions API",
			Triggers:    []string{"interactions api", "gemini agent", "multi-turn", "background research", "chat session"},
			Keywords:    []string{"interactions", "agent", "multimodal", "streaming", "turns"},
			Instructions: `# Gemini Interactions API Agent
Manage stateful agent conversations and background tool execution loops.`,
		},
		{
			ID:          "gemini_live_audio_streamer",
			Name:        "Gemini Live API & Audio Streamer",
			Category:    "AI & Google Cloud",
			Description: "Connects to Gemini Live API via WebSockets for real-time bidirectional audio/video streaming and voice activity detection",
			ToolName:    "gemini_live_stream",
			ToolDesc:    "Manage real-time WebSocket streaming with Gemini Live API",
			Triggers:    []string{"gemini live", "live api", "websocket audio", "bidirectional streaming", "voice activity detection"},
			Keywords:    []string{"live", "audio", "streaming", "websocket", "vad", "voice"},
			Instructions: `# Gemini Live API & Audio Streamer
Build sub-200ms voice-to-voice and streaming video agent integrations.`,
		},
		{
			ID:          "gemini_omni_video_editor",
			Name:        "Gemini Omni Video & Multimodal Editor",
			Category:    "AI & Google Cloud",
			Description: "Performs generative video editing, frame-to-video transitions, ffmpeg preprocessing, and synchronized audio reconstruction",
			ToolName:    "gemini_omni_edit",
			ToolDesc:    "Execute generative video editing with Gemini Omni Flash",
			Triggers:    []string{"omni flash", "video editing", "generative video", "text to video", "ffmpeg video"},
			Keywords:    []string{"video", "omni", "multimodal", "ffmpeg", "transition", "animation"},
			Instructions: `# Gemini Omni Video & Multimodal Editor
Process and transform video sequences with generative frame interpolation and sound synthesis.`,
		},
		{
			ID:          "bigquery_sql_optimizer",
			Name:        "BigQuery SQL Optimizer & Slot Tuner",
			Category:    "Data & Google Cloud",
			Description: "Optimizes BigQuery SQL queries, partition pruning, clustering keys, approximate aggregations, and slot consumption",
			ToolName:    "bigquery_sql_tune",
			ToolDesc:    "Tune BigQuery SQL queries for minimum slot-hours and bytes scanned",
			Triggers:    []string{"bigquery sql", "partition pruning", "clustering", "slot reservation", "bytes billed", "bigquery performance"},
			Keywords:    []string{"bigquery", "sql", "tuning", "partitioning", "slots", "gcp"},
			Instructions: `# BigQuery SQL Optimizer & Slot Tuner
Eliminate full table scans and reduce BigQuery billing costs with partition pruning.`,
		},
		{
			ID:          "bigquery_ai_ml_forecaster",
			Name:        "BigQuery ML & Time Series Forecaster",
			Category:    "Data & Google Cloud",
			Description: "Builds BigQuery ML ARIMA_PLUS models, autoencoder anomaly detection, and integrates remote Gemini LLM models in SQL",
			ToolName:    "bigquery_ml_forecast",
			ToolDesc:    "Train and evaluate BigQuery ML models and time-series forecasts",
			Triggers:    []string{"bigquery ml", "arima_plus", "time series forecast", "bqml", "anomaly detection"},
			Keywords:    []string{"bqml", "ml", "forecast", "arima", "bigquery", "machine learning"},
			Instructions: `# BigQuery ML & Time Series Forecaster
Run in-database machine learning and forecasting natively in BigQuery SQL.`,
		},
		{
			ID:          "bigquery_graph_gql_engine",
			Name:        "BigQuery Graph & GQL Engine",
			Category:    "Data & Google Cloud",
			Description: "Creates property graphs and executes Graph Query Language (GQL) path traversals directly in BigQuery",
			ToolName:    "bigquery_gql_query",
			ToolDesc:    "Query property graphs and topologies in BigQuery using GQL",
			Triggers:    []string{"bigquery graph", "gql", "property graph", "graph query language", "graph traversal"},
			Keywords:    []string{"graph", "gql", "bigquery", "nodes", "edges", "path"},
			Instructions: `# BigQuery Graph & GQL Engine
Model entity relationships and graph topologies using ISO GQL standard queries in BigQuery.`,
		},
		{
			ID:          "bigquery_bigframes_pandas",
			Name:        "BigQuery DataFrames & BigFrames",
			Category:    "Data & Google Cloud",
			Description: "Executes pandas and scikit-learn DataFrame operations on petabyte datasets backed by BigQuery compute engine",
			ToolName:    "bigframes_dataframe_ops",
			ToolDesc:    "Execute BigFrames pandas operations on BigQuery tables",
			Triggers:    []string{"bigframes", "bigquery dataframe", "pandas bigquery", "scikit-learn bigquery"},
			Keywords:    []string{"bigframes", "pandas", "dataframe", "python", "bigquery"},
			Instructions: `# BigQuery DataFrames & BigFrames
Process massive datasets using familiar pandas APIs without pulling data to local memory.`,
		},
		{
			ID:          "dataform_sqlx_pipeline",
			Name:        "Dataform SQLX Pipeline Engineer",
			Category:    "Data & Google Cloud",
			Description: "Authors Dataform SQLX pipeline definitions, assertions, incremental table materializations, and workflow_settings.yaml",
			ToolName:    "dataform_compile_run",
			ToolDesc:    "Compile and execute Dataform SQLX transformation pipelines",
			Triggers:    []string{"dataform", "sqlx", "dataform pipeline", "assertions", "workflow_settings"},
			Keywords:    []string{"dataform", "sqlx", "elt", "bigquery", "transformation"},
			Instructions: `# Dataform SQLX Pipeline Engineer
Author enterprise ELT data pipelines in Google Cloud Dataform.`,
		},
		{
			ID:          "dbt_bigquery_analytics",
			Name:        "dbt BigQuery Analytics Engineer",
			Category:    "Data & Analytics",
			Description: "Develops modular dbt models, Jinja macros, generic schema tests, incremental materializations, and documentation DAGs",
			ToolName:    "dbt_bigquery_run",
			ToolDesc:    "Execute and test dbt models against BigQuery data warehouses",
			Triggers:    []string{"dbt", "dbt-bigquery", "dbt run", "dbt test", "jinja macro", "analytics engineering"},
			Keywords:    []string{"dbt", "analytics", "jinja", "models", "data warehouse"},
			Instructions: `# dbt BigQuery Analytics Engineer
Build reliable, tested transformation DAGs with dbt and BigQuery.`,
		},
		{
			ID:          "gcp_dataflow_beam_runner",
			Name:        "Apache Beam & Dataflow Runner",
			Category:    "Data & Google Cloud",
			Description: "Authors Apache Beam Java/Python pipelines, packages Flex Templates, and troubleshoots streaming autoscaling bottlenecks",
			ToolName:    "dataflow_job_manager",
			ToolDesc:    "Deploy and inspect Apache Beam pipelines on Google Cloud Dataflow",
			Triggers:    []string{"dataflow", "apache beam", "beam pipeline", "flex template", "streaming job"},
			Keywords:    []string{"dataflow", "beam", "streaming", "batch", "gcp"},
			Instructions: `# Apache Beam & Dataflow Runner
Build and manage real-time streaming data processing pipelines on Dataflow.`,
		},
		{
			ID:          "gcp_spark_dataproc_serverless",
			Name:        "Dataproc Serverless & Apache Spark",
			Category:    "Data & Google Cloud",
			Description: "Submits PySpark and Spark SQL batch jobs, manages Dataproc Serverless sessions, and queries BigLake Iceberg tables",
			ToolName:    "dataproc_spark_submit",
			ToolDesc:    "Submit and monitor PySpark batch jobs on Dataproc Serverless",
			Triggers:    []string{"dataproc", "spark", "pyspark", "spark sql", "iceberg catalog", "dataproc serverless"},
			Keywords:    []string{"spark", "dataproc", "pyspark", "iceberg", "biglake", "serverless"},
			Instructions: `# Dataproc Serverless & Apache Spark
Run scalable Apache Spark compute jobs with zero infrastructure management.`,
		},
		{
			ID:          "gcp_composer_airflow_orchestrator",
			Name:        "Cloud Composer & Airflow Orchestrator",
			Category:    "Data & Google Cloud",
			Description: "Orchestrates Apache Airflow DAGs on Cloud Composer (MSAA Gen 2 & 3), manages task dependencies, and debugs failures",
			ToolName:    "composer_dag_trigger",
			ToolDesc:    "Trigger and monitor Apache Airflow DAGs in Cloud Composer",
			Triggers:    []string{"composer", "airflow", "dag", "cloud composer", "msaa", "task dependency"},
			Keywords:    []string{"composer", "airflow", "orchestration", "dag", "gcp"},
			Instructions: `# Cloud Composer & Airflow Orchestrator
Manage complex multi-stage data orchestration workflows with Apache Airflow on GCP.`,
		},
		{
			ID:          "gcs_security_saif_auditor",
			Name:        "Cloud Storage Security & SAIF Auditor",
			Category:    "Security & Google Cloud",
			Description: "Assesses Google Cloud Storage bucket IAM permissions, public access prevention, SAIF compliance, and KMS CMEK keys",
			ToolName:    "gcs_security_audit",
			ToolDesc:    "Audit Google Cloud Storage buckets for SAIF compliance and security risks",
			Triggers:    []string{"gcs security", "saif compliance", "bucket iam", "public access prevention", "cmek", "gcs audit"},
			Keywords:    []string{"gcs", "security", "saif", "cloud storage", "bucket", "iam"},
			Instructions: `# Cloud Storage Security & SAIF Auditor
Verify that cloud storage assets comply with Google Secure AI Framework (SAIF) standards.`,
		},
		{
			ID:          "gcp_data_assets_discovery",
			Name:        "Dataplex & GCP Data Asset Discoverer",
			Category:    "Data & Google Cloud",
			Description: "Searches enterprise data catalogs across BigQuery, BigLake, Spanner, and Dataplex with schema and metadata profiling",
			ToolName:    "dataplex_asset_discover",
			ToolDesc:    "Discover data assets, tables, and schemas in Dataplex and BigQuery",
			Triggers:    []string{"dataplex", "data assets", "discover tables", "data catalog", "spanner table", "table schema", "bigquery datasets", "discover data assets", "gcp data assets"},
			Keywords:    []string{"dataplex", "catalog", "assets", "schema", "discovery", "bigquery", "datasets"},
			Instructions: `# Dataplex & GCP Data Asset Discoverer
Locate datasets, schemas, and table governance policies across cloud data assets.`,
		},
		{
			ID:          "building_data_apps_dashboard",
			Name:        "Data App & Dashboard Builder",
			Category:    "Frontend & Apps",
			Description: "Scaffolds interactive data visualization web apps using React + Vite or Streamlit with Gemini Data Analytics chat",
			ToolName:    "build_data_app_scaffold",
			ToolDesc:    "Scaffold React or Streamlit data applications with BigQuery integration",
			Triggers:    []string{"data app", "dashboard", "streamlit", "react vite", "data visualization", "chat with your data"},
			Keywords:    []string{"dashboard", "streamlit", "react", "frontend", "charts", "analytics"},
			Instructions: `# Data App & Dashboard Builder
Build interactive dashboards that connect directly to BigQuery and provide conversational analytics.`,
		},
		{
			ID:          "data_autocleaning_normalizer",
			Name:        "Automated Data Quality & Normalizer",
			Category:    "Data & Analytics",
			Description: "Applies automated data cleansing, schema mapping, null value imputation, duplicate removal, and type normalization",
			ToolName:    "data_clean_transform",
			ToolDesc:    "Execute automated data cleansing and schema normalization",
			Triggers:    []string{"data cleaning", "schema mapping", "data quality", "null imputation", "deduplication", "data autocleaning"},
			Keywords:    []string{"cleaning", "quality", "normalizer", "schema", "etl", "imputation"},
			Instructions: `# Automated Data Quality & Normalizer
Sanitize raw incoming data with standard type casting and automated anomaly handling.`,
		},

		// === 3. Life Sciences, Genomics & Healthcare Skills ===
		{
			ID:          "alphafold_structure_predictor",
			Name:        "AlphaFold Structure & Confidence Analyzer",
			Category:    "Life Sciences",
			Description: "Retrieves AlphaFold 3D protein structure predictions, per-residue pLDDT confidence scores, and domain boundaries",
			ToolName:    "alphafold_fetch_analyze",
			ToolDesc:    "Fetch and analyze AlphaFold 3D protein structures and pLDDT scores",
			Triggers:    []string{"alphafold", "uniprot structure", "plddt", "protein structure", "domain boundary", "alphafold database"},
			Keywords:    []string{"alphafold", "protein", "plddt", "structure", "uniprot", "biology"},
			Instructions: `# AlphaFold Structure & Confidence Analyzer
Query predicted 3D protein coordinates and assess structural disorder regions.`,
		},
		{
			ID:          "pubmed_clinical_search",
			Name:        "PubMed Biomedical Literature Search",
			Category:    "Life Sciences",
			Description: "Searches PubMed NCBI biomedical publications, fetches abstracts, clinical trial citations, and MeSH indexing",
			ToolName:    "pubmed_search_articles",
			ToolDesc:    "Search PubMed database for peer-reviewed biomedical literature",
			Triggers:    []string{"pubmed", "ncbi", "pmid", "biomedical literature", "clinical trial paper", "mesh terms"},
			Keywords:    []string{"pubmed", "literature", "journal", "mesh", "medicine", "biology"},
			Instructions: `# PubMed Biomedical Literature Search
Discover medical papers, clinical studies, and biological mechanisms from PubMed.`,
		},
		{
			ID:          "clinical_trials_gov_api",
			Name:        "ClinicalTrials.gov Protocol Matcher",
			Category:    "Life Sciences",
			Description: "Queries ClinicalTrials.gov APIv2 for interventional trials, disease conditions, NCT identifiers, and eligibility criteria",
			ToolName:    "clinical_trials_lookup",
			ToolDesc:    "Search ClinicalTrials.gov for clinical study protocols and patient matching",
			Triggers:    []string{"clinical trials", "clinicaltrials.gov", "nct", "trial eligibility", "phase 3 trial", "interventional study"},
			Keywords:    []string{"trials", "nct", "clinical", "pharma", "recruiting", "study"},
			Instructions: `# ClinicalTrials.gov Protocol Matcher
Inspect active clinical trials, endpoint outcomes, and sponsor portfolios.`,
		},
		{
			ID:          "openfda_regulatory_auditor",
			Name:        "openFDA Regulatory & Safety Auditor",
			Category:    "Life Sciences",
			Description: "Queries openFDA API for adverse drug event reports, 510(k) medical device clearances, NDC numbers, and recalls",
			ToolName:    "openfda_query_records",
			ToolDesc:    "Query openFDA database for drug safety, recalls, and device approvals",
			Triggers:    []string{"openfda", "fda adverse event", "510k", "ndc lookup", "drug recall", "fda approval"},
			Keywords:    []string{"fda", "openfda", "safety", "adverse", "pharma", "clearance"},
			Instructions: `# openFDA Regulatory & Safety Auditor
Audit FDA regulatory filings, drug interaction signals, and medical device clearances.`,
		},
		{
			ID:          "pubchem_cheminformatics",
			Name:        "PubChem Chemical Compound Database",
			Category:    "Life Sciences",
			Description: "Resolves chemical names, SMILES strings, PubChem CIDs, molecular weights, 2D/3D structures, and bioassays",
			ToolName:    "pubchem_compound_search",
			ToolDesc:    "Search PubChem for chemical compounds, SMILES, and bioactivity data",
			Triggers:    []string{"pubchem", "smiles", "cid", "chemical structure", "cheminformatics", "molecular formula"},
			Keywords:    []string{"pubchem", "chemistry", "smiles", "compound", "molecule", "bioassay"},
			Instructions: `# PubChem Chemical Compound Database
Look up chemical entities, molecular properties, and pharmacological assays.`,
		},
		{
			ID:          "reactome_pathway_analyzer",
			Name:        "Reactome Pathway Enrichment Analyzer",
			Category:    "Life Sciences",
			Description: "Performs biological pathway enrichment, reaction inputs/outputs, topological hierarchy, and diagram export",
			ToolName:    "reactome_analyze_pathway",
			ToolDesc:    "Analyze biological pathways and gene list enrichment with Reactome",
			Triggers:    []string{"reactome", "pathway enrichment", "biological pathway", "reaction participants", "cellular pathway"},
			Keywords:    []string{"reactome", "pathway", "signaling", "enrichment", "reaction", "genes"},
			Instructions: `# Reactome Pathway Enrichment Analyzer
Map sets of genes to annotated metabolic and signaling biological pathways.`,
		},
		{
			ID:          "string_ppi_network_analyzer",
			Name:        "STRING Protein Interaction Network",
			Category:    "Life Sciences",
			Description: "Queries STRING database for physical and functional protein-protein interaction networks and confidence scores",
			ToolName:    "string_ppi_query",
			ToolDesc:    "Retrieve protein-protein interaction networks and confidence scores from STRING",
			Triggers:    []string{"string database", "protein interaction", "ppi network", "protein partners", "interactome"},
			Keywords:    []string{"string", "ppi", "interactome", "protein", "network", "confidence"},
			Instructions: `# STRING Protein Interaction Network
Explore molecular interaction partners and functional protein modules.`,
		},
		{
			ID:          "uniprot_kb_annotator",
			Name:        "UniProt Knowledgebase Annotator",
			Category:    "Life Sciences",
			Description: "Fetches UniProtKB protein sequence records, active catalytic residues, isoforms, subcellular locations, and GO tags",
			ToolName:    "uniprot_fetch_protein",
			ToolDesc:    "Look up protein metadata and sequences in UniProtKB",
			Triggers:    []string{"uniprot", "uniprotkb", "protein sequence", "catalytic site", "isoform", "protein annotation"},
			Keywords:    []string{"uniprot", "protein", "sequence", "fasta", "annotation", "swissprot"},
			Instructions: `# UniProt Knowledgebase Annotator
Retrieve functional annotations and amino acid sequences from Swiss-Prot/UniProtKB.`,
		},
		{
			ID:          "chembl_target_bioactivity",
			Name:        "ChEMBL Drug Target & Bioactivity Database",
			Category:    "Life Sciences",
			Description: "Queries ChEMBL for bioactive small molecules, IC50/Ki target affinities, approved drug mechanisms, and clinical phase",
			ToolName:    "chembl_target_query",
			ToolDesc:    "Query ChEMBL for drug targets, IC50/Ki bioactivities, and approved drugs",
			Triggers:    []string{"chembl", "ic50", "ki affinity", "drug target", "bioactivity", "target inhibition"},
			Keywords:    []string{"chembl", "pharma", "ic50", "target", "drug", "affinity"},
			Instructions: `# ChEMBL Drug Target & Bioactivity Database
Evaluate compound binding affinities and target selectivity profiles.`,
		},
		{
			ID:          "arxiv_paper_search_fetcher",
			Name:        "ArXiv Scientific Preprints Search",
			Category:    "Research & Literature",
			Description: "Searches arXiv repository across Computer Science, Quantitative Biology, Physics, and Math, extracting full text",
			ToolName:    "arxiv_search_papers",
			ToolDesc:    "Search arXiv preprints and download paper metadata and text",
			Triggers:    []string{"arxiv", "preprint", "arxiv id", "research paper", "arxiv search", "cs paper"},
			Keywords:    []string{"arxiv", "paper", "research", "academic", "preprint", "science"},
			Instructions: `# ArXiv Scientific Preprints Search
Search and review cutting-edge preprints across scientific fields.`,
		},
		{
			ID:          "openalex_scholarly_graph",
			Name:        "OpenAlex Global Scholarly Graph",
			Category:    "Research & Literature",
			Description: "Navigates 250M+ scholarly works, author h-index bibliometrics, institutional citations, and Open Access DOIs",
			ToolName:    "openalex_query_graph",
			ToolDesc:    "Query the OpenAlex scholarly knowledge graph and bibliometrics",
			Triggers:    []string{"openalex", "scholarly graph", "author h-index", "citation count", "doi lookup", "bibliometrics"},
			Keywords:    []string{"openalex", "scholar", "citations", "hindex", "doi", "publications"},
			Instructions: `# OpenAlex Global Scholarly Graph
Analyze academic research impact and citation topologies.`,
		},
		{
			ID:          "europepmc_fulltext_xml",
			Name:        "Europe PMC Full-Text & Bio-Entities",
			Category:    "Research & Literature",
			Description: "Retrieves open-access biomedical literature in full-text XML and identifies mined bio-entities and chemicals",
			ToolName:    "europepmc_fetch_fulltext",
			ToolDesc:    "Fetch full-text XML articles and entity annotations from Europe PMC",
			Triggers:    []string{"europe pmc", "europepmc", "full text xml", "pmcid", "biomedical fulltext"},
			Keywords:    []string{"europepmc", "pmc", "xml", "fulltext", "biomedical"},
			Instructions: `# Europe PMC Full-Text & Bio-Entities
Access complete scientific literature text and automated text-mining annotations.`,
		},
		{
			ID:          "clinvar_pathogenicity_auditor",
			Name:        "ClinVar Genomic Pathogenicity Auditor",
			Category:    "Genomics",
			Description: "Evaluates human genetic variants against ClinVar for Pathogenic, Likely Pathogenic, Benign, or VUS classifications",
			ToolName:    "clinvar_variant_lookup",
			ToolDesc:    "Look up human genomic variant pathogenicity and clinical evidence in ClinVar",
			Triggers:    []string{"clinvar", "pathogenicity", "variant clinical significance", "vus", "pathogenic variant", "acmg"},
			Keywords:    []string{"clinvar", "genomics", "pathogenic", "variant", "mutation", "clinical"},
			Instructions: `# ClinVar Genomic Pathogenicity Auditor
Determine clinical significance and ACMG classification for human genetic variants.`,
		},
		{
			ID:          "dbsnp_variant_mapper",
			Name:        "dbSNP Variant & rsID Coordinate Mapper",
			Category:    "Genomics",
			Description: "Maps rsIDs to GRCh38 genomic coordinates, HGVS strings, minor allele frequencies, and indel classifications",
			ToolName:    "dbsnp_lookup_variant",
			ToolDesc:    "Resolve and map genetic variants (SNPs/indels) in NCBI dbSNP",
			Triggers:    []string{"dbsnp", "rsid", "snp", "genomic coordinates", "grch38 variant", "indel lookup"},
			Keywords:    []string{"dbsnp", "rsid", "snp", "genomics", "allele", "ncbi"},
			Instructions: `# dbSNP Variant & rsID Coordinate Mapper
Translate variant identifiers into standard genomic notation and allele frequencies.`,
		},
		{
			ID:          "gnomad_allele_frequency",
			Name:        "gnomAD Population Allele Frequency",
			Category:    "Genomics",
			Description: "Queries the Genome Aggregation Database for population allele frequencies and gene constraint metrics (pLI, LOEUF)",
			ToolName:    "gnomad_query_frequency",
			ToolDesc:    "Query gnomAD for population variant frequencies and gene constraint metrics",
			Triggers:    []string{"gnomad", "allele frequency", "pli constraint", "loeuf", "population genomics", "loss of function"},
			Keywords:    []string{"gnomad", "population", "allele", "frequency", "constraint", "genomics"},
			Instructions: `# gnomAD Population Allele Frequency
Determine the rarity of genetic variants across global ancestral populations.`,
		},
		{
			ID:          "interpro_domain_architect",
			Name:        "InterPro Protein Family & Pfam Architect",
			Category:    "Life Sciences",
			Description: "Scans protein sequences for functional domains, Pfam families, active sites, and deep-learning InterPro-N models",
			ToolName:    "interpro_domain_scan",
			ToolDesc:    "Scan protein sequences against InterPro and Pfam domain signatures",
			Triggers:    []string{"interpro", "pfam", "protein domain", "domain architecture", "interpro-n", "protein family"},
			Keywords:    []string{"interpro", "pfam", "domain", "family", "protein", "signature"},
			Instructions: `# InterPro Protein Family & Pfam Architect
Classify protein families and identify conserved catalytic and binding domains.`,
		},
		{
			ID:          "pdb_macromolecule_structure",
			Name:        "Protein Data Bank (PDB) Macromolecule Structure",
			Category:    "Life Sciences",
			Description: "Downloads experimentally determined 3D atomic coordinates (.cif/.pdb) and binding ligand metadata from the Protein Data Bank",
			ToolName:    "pdb_fetch_structure",
			ToolDesc:    "Fetch 3D crystal structures and experimental metadata from PDB",
			Triggers:    []string{"pdb", "protein data bank", "pdb structure", "crystal structure", "pdb id", "cryo-em structure"},
			Keywords:    []string{"pdb", "structure", "crystal", "macromolecule", "ligand", "cif"},
			Instructions: `# Protein Data Bank (PDB) Macromolecule Structure
Query high-resolution atomic coordinates for proteins and drug complexes.`,
		},
		{
			ID:          "pymol_molecular_visualizer",
			Name:        "PyMOL 3D Molecular Rendering Engine",
			Category:    "Life Sciences",
			Description: "Generates PyMOL rendering scripts, superpositions, binding pocket raytracing, and surface electrostatic views",
			ToolName:    "pymol_render_scene",
			ToolDesc:    "Generate PyMOL visualization scripts and render molecular scenes",
			Triggers:    []string{"pymol", "molecular visualization", "render protein", "binding pocket render", "structural alignment"},
			Keywords:    []string{"pymol", "rendering", "3d", "visualization", "raytracing", "superposition"},
			Instructions: `# PyMOL 3D Molecular Rendering Engine
Produce publication-quality structural renders and binding site alignments.`,
		},
		{
			ID:          "quickgo_ontology_mapper",
			Name:        "QuickGO Gene Ontology Term Navigator",
			Category:    "Life Sciences",
			Description: "Maps Gene Ontology (GO) terms across Biological Process, Molecular Function, and Cellular Component hierarchies",
			ToolName:    "quickgo_fetch_terms",
			ToolDesc:    "Query QuickGO for Gene Ontology terms and ECO evidence codes",
			Triggers:    []string{"quickgo", "gene ontology", "go term", "biological process", "molecular function", "eco code"},
			Keywords:    []string{"quickgo", "go", "ontology", "biological process", "molecular function"},
			Instructions: `# QuickGO Gene Ontology Term Navigator
Navigate the functional taxonomy of genes and physiological processes.`,
		},
		{
			ID:          "jaspar_tfbs_matrix_caller",
			Name:        "JASPAR Transcription Factor Binding Profiles",
			Category:    "Genomics",
			Description: "Queries JASPAR database for Transcription Factor binding Position Weight Matrices (PWM) and MEME formatted motifs",
			ToolName:    "jaspar_matrix_query",
			ToolDesc:    "Fetch transcription factor binding profiles and PWM matrices from JASPAR",
			Triggers:    []string{"jaspar", "transcription factor", "tfbs", "pwm matrix", "binding motif", "meme format"},
			Keywords:    []string{"jaspar", "transcription", "tfbs", "pwm", "motif", "dna binding"},
			Instructions: `# JASPAR Transcription Factor Binding Profiles
Retrieve sequence binding motifs for regulatory proteins and transcription factors.`,
		},
		{
			ID:          "encode_screen_ccre_query",
			Name:        "ENCODE SCREEN cis-Regulatory Elements",
			Category:    "Genomics",
			Description: "Searches ENCODE Registry of candidate cis-Regulatory Elements (cCREs) across human cell lines and ChIP-seq peaks",
			ToolName:    "encode_ccre_search",
			ToolDesc:    "Search candidate cis-regulatory elements in ENCODE SCREEN",
			Triggers:    []string{"encode", "ccre", "screen", "regulatory element", "promoter enhancer", "chip-seq peak"},
			Keywords:    []string{"encode", "ccre", "screen", "enhancer", "promoter", "epigenomics"},
			Instructions: `# ENCODE SCREEN cis-Regulatory Elements
Inspect genomic enhancers, promoters, and insulator elements active in specific tissues.`,
		},
		{
			ID:          "predicting_the_past_epigraphy",
			Name:        "Ancient Epigraphy & Text Restorer",
			Category:    "Humanities & AI",
			Description: "Restores missing characters, dates, and attributes ancient Greek and Latin inscriptions using Aeneas and Ithaca models",
			ToolName:    "epigraphy_restore_text",
			ToolDesc:    "Restore, date, and locate ancient epigraphic texts using Aeneas/Ithaca",
			Triggers:    []string{"epigraphy", "ancient greek", "latin inscription", "ithaca", "aeneas", "text restoration"},
			Keywords:    []string{"epigraphy", "ancient", "greek", "latin", "history", "restoration"},
			Instructions: `# Ancient Epigraphy & Text Restorer
Reconstruct damaged historical texts and predict geographical provenance.`,
		},
		{
			ID:          "alphagenome_variant_impact",
			Name:        "AlphaGenome Non-Coding Regulatory Predictor",
			Category:    "Genomics",
			Description: "Predicts non-coding variant impact on gene expression, chromatin accessibility (DNase), and histone marks",
			ToolName:    "alphagenome_variant_predict",
			ToolDesc:    "Predict regulatory and expression effects of non-coding variants with AlphaGenome",
			Triggers:    []string{"alphagenome", "non-coding variant", "chromatin accessibility", "dnase", "regulatory variant effect"},
			Keywords:    []string{"alphagenome", "genomics", "expression", "noncoding", "chromatin"},
			Instructions: `# AlphaGenome Non-Coding Regulatory Predictor
Assess the functional pathogenicity of non-coding genomic mutations.`,
		},

		// === 4. FinTech, Blockchain & SEC Financial Modeling ===
		{
			ID:          "sec_edgar_analyst",
			Name:        "SEC EDGAR Financial Analyst",
			Category:    "Finance",
			Description: "Retrieves 10-K, 10-Q, 8-K filings and exact GAAP metrics from SEC EDGAR (eval mock)",
			ToolName:    "sec_edgar_query",
			ToolDesc:    "Query SEC EDGAR database for verified financial metrics",
			Triggers:    []string{"sec", "10-k", "10-q", "gaap", "operating margin", "edgar"},
			Keywords:    []string{"sec", "edgar", "10k", "10q", "gaap", "revenue", "filings"},
			Instructions: `# SEC EDGAR Financial Analyst
Extract certified corporate financial figures directly from audited SEC filings.`,
		},
		{
			ID:          "val_dcf_modeler",
			Name:        "DCF & Valuation Financial Modeler",
			Category:    "Finance",
			Description: "Calculates discounted cash flows, terminal value, and WACC hurdle rates with sensitivity tables",
			ToolName:    "valuation_dcf_calc",
			ToolDesc:    "Compute DCF financial models and valuation",
			Triggers:    []string{"dcf", "valuation", "wacc", "cash flow", "terminal value"},
			Keywords:    []string{"dcf", "valuation", "wacc", "discounted", "finance", "cash flow"},
			Instructions: `# DCF & Valuation Financial Modeler
Model enterprise valuation and cash flow projections under multiple discount rate scenarios.`,
		},
		{
			ID:          "financesec_lattice_mcp",
			Name:        "Calera FinanceSec Volumetric Lattice",
			Category:    "Finance",
			Description: "Queries Calera FinanceSec lattice for 12 valuation packs, FOMC dot plots, and Treasury yield curves (eval mock)",
			ToolName:    "finsec_lattice_query",
			ToolDesc:    "Query FinanceSec certified volumetric financial lattice",
			Triggers:    []string{"finsec", "financesec", "fomc dot plot", "treasury yield curve", "valuation pack", "merkle audit"},
			Keywords:    []string{"finsec", "calera", "lattice", "treasury", "fomc", "valuation"},
			Instructions: `# Calera FinanceSec Volumetric Lattice
Query zero-hallucination verified financial metrics and macroeconomic dot plots.`,
		},
		{
			ID:          "stripe_billing_ops",
			Name:        "Stripe Billing & Payment Reconciler",
			Category:    "FinTech",
			Description: "Reconciles customer invoices, subscription billing, payment intents, and webhooks",
			ToolName:    "stripe_reconciler",
			ToolDesc:    "Reconcile Stripe invoices and webhooks",
			Triggers:    []string{"stripe", "invoice", "payment", "webhook", "billing"},
			Keywords:    []string{"stripe", "billing", "invoice", "payment intent", "subscription"},
			Instructions: `# Stripe Billing & Payment Reconciler
Process and reconcile payment flows, subscription upgrades, and webhook events.`,
		},
		{
			ID:          "plaid_bank_account_link",
			Name:        "Plaid Open Banking & ACH Verifier",
			Category:    "FinTech",
			Description: "Verifies bank account routing, real-time balances, ACH identity tokens, and transaction history via Plaid",
			ToolName:    "plaid_verify_account",
			ToolDesc:    "Verify bank accounts and check real-time balances via Plaid",
			Triggers:    []string{"plaid", "bank verification", "ach routing", "bank balance", "open banking"},
			Keywords:    []string{"plaid", "banking", "ach", "balance", "transactions", "fintech"},
			Instructions: `# Plaid Open Banking & ACH Verifier
Authenticate bank accounts and verify liquidity before initiating transfers.`,
		},
		{
			ID:          "alpaca_market_order_broker",
			Name:        "Alpaca Market Order Broker",
			Category:    "FinTech",
			Description: "Executes equity and options orders, inspects portfolio buying power, and sets bracket stop-loss orders on Alpaca",
			ToolName:    "alpaca_place_order",
			ToolDesc:    "Submit and manage equity market and limit orders via Alpaca",
			Triggers:    []string{"alpaca", "market order", "limit order", "bracket order", "buying power", "stock trade"},
			Keywords:    []string{"alpaca", "trading", "stocks", "orders", "portfolio", "broker"},
			Instructions: `# Alpaca Market Order Broker
Execute automated trading algorithms within strict portfolio risk limits.`,
		},
		{
			ID:          "bloomberg_bpipe_feed",
			Name:        "Bloomberg B-PIPE Market Feed",
			Category:    "Finance",
			Description: "Streams institutional Level 2 market data, tick depth, swap rates, and yield curve spreads from Bloomberg B-PIPE",
			ToolName:    "bloomberg_market_quote",
			ToolDesc:    "Fetch real-time market depth and pricing from Bloomberg B-PIPE",
			Triggers:    []string{"bloomberg", "bpipe", "market depth", "level 2 quote", "yield spread", "orderbook"},
			Keywords:    []string{"bloomberg", "market", "quote", "orderbook", "bpipe", "ticks"},
			Instructions: `# Bloomberg B-PIPE Market Feed
Access institutional real-time pricing and liquidity depth.`,
		},

		// === 5. Cloud, Infrastructure, DevOps & Orchestration ===
		{
			ID:          "docker_container_ops",
			Name:        "Docker Container Lifecycle Manager",
			Category:    "DevOps",
			Description: "Manages docker containers, inspects container logs, and restarts microservices",
			ToolName:    "docker_manager",
			ToolDesc:    "Inspect and manage Docker container lifecycle",
			Triggers:    []string{"docker", "container", "restart", "container health", "microservice"},
			Keywords:    []string{"docker", "container", "image", "daemon", "dockerfile"},
			Instructions: `# Docker Container Lifecycle Manager
Manage containerized applications and examine runtime container health.`,
		},
		{
			ID:          "kubernetes_cluster_ops",
			Name:        "Kubernetes Cluster Orchestrator",
			Category:    "DevOps",
			Description: "Scales pods, updates deployments, and manages namespaces in Kubernetes clusters",
			ToolName:    "kubectl_orchestrator",
			ToolDesc:    "Execute kubectl operations on Kubernetes clusters",
			Triggers:    []string{"kubernetes", "k8s", "kubectl", "replicas", "scale", "namespace", "deployment"},
			Keywords:    []string{"k8s", "kubernetes", "pod", "deployment", "kubectl", "cluster"},
			Instructions: `# Kubernetes Cluster Orchestrator
Automate cluster deployments, rolling updates, and autoscaling policies.`,
		},
		{
			ID:          "helm_chart_deployer",
			Name:        "Helm Kubernetes Chart Deployer",
			Category:    "DevOps",
			Description: "Renders Helm templates and installs chart releases into Kubernetes",
			ToolName:    "helm_install_release",
			ToolDesc:    "Install or upgrade Helm chart releases",
			Triggers:    []string{"helm", "chart", "kubernetes deploy", "release", "helm install"},
			Keywords:    []string{"helm", "chart", "values", "release", "k8s"},
			Instructions: `# Helm Kubernetes Chart Deployer
Package and deploy parameterized application manifests with Helm.`,
		},
		{
			ID:          "terraform_infra_plan",
			Name:        "Terraform Infrastructure Planner",
			Category:    "DevOps",
			Description: "Generates and validates Terraform execution plans for cloud infrastructure",
			ToolName:    "terraform_plan",
			ToolDesc:    "Run terraform plan and validate HCL",
			Triggers:    []string{"terraform", "hcl", "infrastructure", "iac", "terraform plan"},
			Keywords:    []string{"terraform", "iac", "hcl", "infrastructure", "plan", "apply"},
			Instructions: `# Terraform Infrastructure Planner
Inspect planned infrastructure modifications and prevent configuration drift.`,
		},
		{
			ID:          "ansible_playbook_runner",
			Name:        "Ansible Automation Playbook Runner",
			Category:    "DevOps",
			Description: "Executes configuration management playbooks across fleet of servers",
			ToolName:    "ansible_run_playbook",
			ToolDesc:    "Execute Ansible configuration playbooks",
			Triggers:    []string{"ansible", "playbook", "automation", "sysadmin", "inventory"},
			Keywords:    []string{"ansible", "playbook", "yaml", "ssh", "automation"},
			Instructions: `# Ansible Automation Playbook Runner
Enforce idempotent configuration states across cloud server fleets.`,
		},
		{
			ID:          "argo_workflows_cicd",
			Name:        "Argo Workflows Cloud-Native Engine",
			Category:    "DevOps",
			Description: "Orchestrates DAG-based container workflows on Kubernetes clusters",
			ToolName:    "argo_submit_workflow",
			ToolDesc:    "Submit and monitor Argo Workflow DAGs",
			Triggers:    []string{"argo", "workflows", "dag", "kubernetes workflow", "argo submit"},
			Keywords:    []string{"argo", "workflows", "dag", "kubernetes", "ci"},
			Instructions: `# Argo Workflows Cloud-Native Engine
Execute containerized batch computation graphs on Kubernetes.`,
		},
		{
			ID:          "github_actions_ci",
			Name:        "GitHub Actions CI/CD Pipeline Manager",
			Category:    "DevOps",
			Description: "Triggers workflow dispatches and monitors CI test runners",
			ToolName:    "github_actions_trigger",
			ToolDesc:    "Trigger and monitor GitHub Actions workflows",
			Triggers:    []string{"github actions", "ci", "cd", "pipeline", "workflow dispatch"},
			Keywords:    []string{"github", "actions", "ci", "cd", "workflow", "runner"},
			Instructions: `# GitHub Actions CI/CD Pipeline Manager
Trigger continuous integration jobs and monitor build artifacts.`,
		},
		{
			ID:          "github_pr_auditor",
			Name:        "GitHub Pull Request Reviewer",
			Category:    "DevOps",
			Description: "Audits PR changes, runs static analysis, and submits review comments",
			ToolName:    "github_pr_review",
			ToolDesc:    "Review and audit GitHub pull requests",
			Triggers:    []string{"github", "pr", "pull request", "review", "pr diff"},
			Keywords:    []string{"github", "pr", "review", "diff", "pull request"},
			Instructions: `# GitHub Pull Request Reviewer
Review branch diffs and provide automated inline feedback.`,
		},
		{
			ID:          "git_code_patcher",
			Name:        "Git & AST Code Patcher",
			Category:    "DevOps",
			Description: "Inspects repository AST symbol trees and generates unified diff git patches",
			ToolName:    "git_diff_patcher",
			ToolDesc:    "Synthesize git diff patch for code refactoring",
			Triggers:    []string{"git", "diff", "patch", "function", "ast", "refactor"},
			Keywords:    []string{"git", "patch", "diff", "refactor", "ast", "code"},
			Instructions: `# Git & AST Code Patcher
Generate clean, collision-free unified diffs for automated code modifications.`,
		},
		{
			ID:          "aws_s3_storage",
			Name:        "AWS S3 Cloud Storage Manager",
			Category:    "Cloud",
			Description: "Uploads, downloads, and manages bucket lifecycle policies on AWS S3",
			ToolName:    "s3_bucket_ops",
			ToolDesc:    "Manage AWS S3 storage buckets and objects",
			Triggers:    []string{"s3", "aws", "bucket", "storage", "aws s3"},
			Keywords:    []string{"s3", "aws", "bucket", "object", "storage"},
			Instructions: `# AWS S3 Cloud Storage Manager
Manage cloud object storage and presigned access authorizations.`,
		},
		{
			ID:          "aws_lambda_serverless",
			Name:        "AWS Lambda Serverless Function Manager",
			Category:    "Cloud",
			Description: "Deploys Python/Go functions and configures EventBridge trigger rules",
			ToolName:    "aws_lambda_invoke",
			ToolDesc:    "Invoke and deploy AWS Lambda serverless functions",
			Triggers:    []string{"aws lambda", "serverless", "eventbridge", "lambda function"},
			Keywords:    []string{"lambda", "aws", "serverless", "eventbridge", "functions"},
			Instructions: `# AWS Lambda Serverless Function Manager
Deploy event-driven serverless compute handlers.`,
		},
		{
			ID:          "cloudflare_dns_waf",
			Name:        "Cloudflare DNS & WAF Rule Manager",
			Category:    "Networking",
			Description: "Updates DNS A/CNAME records and configures Cloudflare WAF firewall rules",
			ToolName:    "cloudflare_update_dns",
			ToolDesc:    "Update Cloudflare DNS records and WAF",
			Triggers:    []string{"cloudflare", "dns", "waf", "firewall", "dns record"},
			Keywords:    []string{"cloudflare", "dns", "waf", "firewall", "cdn"},
			Instructions: `# Cloudflare DNS & WAF Rule Manager
Manage edge routing and firewall protection policies.`,
		},
		{
			ID:          "nginx_ingress_proxy",
			Name:        "NGINX Ingress & Reverse Proxy",
			Category:    "Networking",
			Description: "Configures upstream routing, SSL termination, and rate limiting rules",
			ToolName:    "nginx_reload_config",
			ToolDesc:    "Configure and reload NGINX reverse proxy",
			Triggers:    []string{"nginx", "proxy", "ssl", "reverse proxy", "ingress"},
			Keywords:    []string{"nginx", "proxy", "reverse", "ssl", "upstream"},
			Instructions: `# NGINX Ingress & Reverse Proxy
Configure high-performance reverse proxy routing and TLS termination.`,
		},

		// === 6. Modern Databases, Search & Message Queues ===
		{
			ID:          "postgres_db_admin",
			Name:        "PostgreSQL Database Administrator",
			Category:    "Database",
			Description: "Executes parameterized SQL queries, schema migrations, and transaction management",
			ToolName:    "postgres_executor",
			ToolDesc:    "Execute SQL queries against PostgreSQL database",
			Triggers:    []string{"postgres", "sql", "database", "table", "transactions table", "settlement status", "committed", "sql query"},
			Keywords:    []string{"postgres", "postgresql", "sql", "database", "query"},
			Instructions: `# PostgreSQL Database Administrator
Execute ACID transactions, partition maintenance, and schema migrations.`,
		},
		{
			ID:          "redis_cache_manager",
			Name:        "Redis In-Memory Cache Manager",
			Category:    "Database",
			Description: "Manages Redis keys, sets TTL, flushes cache keys, and inspects cache eviction",
			ToolName:    "redis_cache_ops",
			ToolDesc:    "Execute operations on Redis in-memory cache",
			Triggers:    []string{"redis", "cache", "ttl", "flush", "orderbook", "redis key"},
			Keywords:    []string{"redis", "cache", "ttl", "in-memory", "key-value"},
			Instructions: `# Redis In-Memory Cache Manager
Optimize caching layers and manage real-time pub/sub channels.`,
		},
		{
			ID:          "mongodb_document_store",
			Name:        "MongoDB NoSQL Document Store",
			Category:    "Database",
			Description: "Runs aggregation pipelines and CRUD operations on MongoDB collections",
			ToolName:    "mongodb_aggregate",
			ToolDesc:    "Execute MongoDB aggregation pipelines",
			Triggers:    []string{"mongodb", "nosql", "document", "mongo", "aggregation pipeline"},
			Keywords:    []string{"mongodb", "mongo", "nosql", "bson", "document"},
			Instructions: `# MongoDB NoSQL Document Store
Run flexible JSON/BSON document aggregation queries.`,
		},
		{
			ID:          "neo4j_graph_cypher",
			Name:        "Neo4j Graph Database & Cypher Engine",
			Category:    "Database",
			Description: "Executes Cypher graph queries to find shortest paths and entity clusters",
			ToolName:    "neo4j_cypher_query",
			ToolDesc:    "Run Cypher queries on Neo4j graph database",
			Triggers:    []string{"neo4j", "cypher", "graph", "nodes", "relationships", "graph traversal"},
			Keywords:    []string{"neo4j", "cypher", "graph", "nodes", "relationships"},
			Instructions: `# Neo4j Graph Database & Cypher Engine
Traverse complex graph topologies and find connected entity clusters.`,
		},
		{
			ID:          "clickhouse_olap_engine",
			Name:        "ClickHouse Real-Time OLAP Engine",
			Category:    "Database",
			Description: "Runs real-time vector and analytical queries on billions of events",
			ToolName:    "clickhouse_sql_query",
			ToolDesc:    "Query ClickHouse columnar database",
			Triggers:    []string{"clickhouse", "olap", "columnar", "realtime", "event analytics"},
			Keywords:    []string{"clickhouse", "olap", "columnar", "analytics", "sql"},
			Instructions: `# ClickHouse Real-Time OLAP Engine
Execute low-latency analytical aggregations over large event streams.`,
		},
		{
			ID:          "snowflake_warehouse_query",
			Name:        "Snowflake Cloud Data Warehouse",
			Category:    "Data",
			Description: "Executes analytics SQL on Snowflake warehouses with auto-scaling",
			ToolName:    "snowflake_query_exec",
			ToolDesc:    "Execute analytical queries on Snowflake",
			Triggers:    []string{"snowflake", "warehouse", "analytics", "data warehouse", "snowpark"},
			Keywords:    []string{"snowflake", "warehouse", "sql", "data", "cloud"},
			Instructions: `# Snowflake Cloud Data Warehouse
Run enterprise data warehouse analytical queries with elastic scaling.`,
		},
		{
			ID:          "elasticsearch_log_search",
			Name:        "Elasticsearch Log Search Engine",
			Category:    "Observability",
			Description: "Searches distributed cluster logs and runs Lucene aggregations",
			ToolName:    "elasticsearch_search",
			ToolDesc:    "Search Elasticsearch cluster indices",
			Triggers:    []string{"elasticsearch", "lucene", "kibana", "cluster logs", "log search"},
			Keywords:    []string{"elasticsearch", "lucene", "logs", "search", "kibana"},
			Instructions: `# Elasticsearch Log Search Engine
Perform full-text BM25 queries and aggregate log events.`,
		},
		{
			ID:          "solr_document_indexer",
			Name:        "Apache Solr Document Indexer",
			Category:    "Search",
			Description: "Indexes text documents and runs full-text BM25 queries",
			ToolName:    "solr_index_ops",
			ToolDesc:    "Manage Solr document collections",
			Triggers:    []string{"solr", "indexer", "search engine", "solr collection"},
			Keywords:    []string{"solr", "search", "indexer", "bm25", "lucene"},
			Instructions: `# Apache Solr Document Indexer
Maintain search indices and faceting collections.`,
		},
		{
			ID:          "sqlite_local_auditor",
			Name:        "SQLite Embedded Ledger Auditor",
			Category:    "Database",
			Description: "Audits local ACID SQLite databases and verifies table indexes",
			ToolName:    "sqlite_verify_integrity",
			ToolDesc:    "Verify SQLite database integrity",
			Triggers:    []string{"sqlite", "embedded", "ledger", "integrity", "sqlite db"},
			Keywords:    []string{"sqlite", "embedded", "acid", "database", "local"},
			Instructions: `# SQLite Embedded Ledger Auditor
Check local embedded database consistency and run fast transactional queries.`,
		},
		{
			ID:          "kafka_stream_processor",
			Name:        "Apache Kafka Stream Processor",
			Category:    "Data",
			Description: "Consumes and produces high-throughput event topics across Kafka brokers",
			ToolName:    "kafka_topic_manager",
			ToolDesc:    "Manage Kafka topics and event streams",
			Triggers:    []string{"kafka", "stream", "topic", "consumer", "producer", "kafka broker"},
			Keywords:    []string{"kafka", "stream", "topic", "event", "broker"},
			Instructions: `# Apache Kafka Stream Processor
Coordinate partitioned event topic ingestion and consumer group commits.`,
		},
		{
			ID:          "rabbitmq_queue_consumer",
			Name:        "RabbitMQ Message Broker Consumer",
			Category:    "Data",
			Description: "Acknowledges messages and manages dead-letter queues in RabbitMQ",
			ToolName:    "rabbitmq_manage_queue",
			ToolDesc:    "Manage RabbitMQ message queues",
			Triggers:    []string{"rabbitmq", "amqp", "queue", "dead letter", "message broker"},
			Keywords:    []string{"rabbitmq", "amqp", "queue", "broker", "routing"},
			Instructions: `# RabbitMQ Message Broker Consumer
Route messages through topic exchanges and manage dead-letter failovers.`,
		},

		// === 7. Vector Databases & AI Search ===
		{
			ID:          "weaviate_vector_search",
			Name:        "Weaviate Vector & Hybrid Search Engine",
			Category:    "Search",
			Description: "Performs ANN vector similarity searches with cosine distance filters",
			ToolName:    "weaviate_vector_query",
			ToolDesc:    "Perform hybrid vector queries on Weaviate",
			Triggers:    []string{"weaviate", "vector", "embedding", "similarity", "hybrid search"},
			Keywords:    []string{"weaviate", "vector", "embedding", "hybrid", "ann"},
			Instructions: `# Weaviate Vector & Hybrid Search Engine
Execute hybrid dense and sparse vector retrieval over multimodal collections.`,
		},
		{
			ID:          "qdrant_vector_store",
			Name:        "Qdrant Vector Database Engine",
			Category:    "Search",
			Description: "Stores dense high-dimensional vectors and runs payload filtered searches",
			ToolName:    "qdrant_search_points",
			ToolDesc:    "Search vectors with payload filters in Qdrant",
			Triggers:    []string{"qdrant", "vector database", "payload filter", "hnsw search"},
			Keywords:    []string{"qdrant", "vector", "hnsw", "payload", "search"},
			Instructions: `# Qdrant Vector Database Engine
Run filtered approximate nearest neighbor searches with payload criteria.`,
		},
		{
			ID:          "milvus_vector_cluster",
			Name:        "Milvus Distributed Vector Database",
			Category:    "Search",
			Description: "Scales billion-scale vector indexes with GPU acceleration",
			ToolName:    "milvus_ann_search",
			ToolDesc:    "Execute ANN vector search on Milvus cluster",
			Triggers:    []string{"milvus", "ann search", "vector index", "billion scale vector"},
			Keywords:    []string{"milvus", "vector", "gpu", "ann", "cluster"},
			Instructions: `# Milvus Distributed Vector Database
Manage enterprise-scale distributed vector indices.`,
		},
		{
			ID:          "pinecone_serverless_index",
			Name:        "Pinecone Serverless Vector Indexer",
			Category:    "Search",
			Description: "Manages Pinecone serverless vector namespaces and executes sub-50ms top-K embeddings retrieval",
			ToolName:    "pinecone_query_vectors",
			ToolDesc:    "Query vectors and metadata in Pinecone serverless indices",
			Triggers:    []string{"pinecone", "pinecone index", "vector namespace", "top-k vectors"},
			Keywords:    []string{"pinecone", "vector", "serverless", "embeddings", "topk"},
			Instructions: `# Pinecone Serverless Vector Indexer
Perform fast vector similarity search across partitioned namespaces.`,
		},
		{
			ID:          "chromadb_local_embeddings",
			Name:        "ChromaDB Embedded Vector Store",
			Category:    "Search",
			Description: "Stores, persists, and queries local collection embeddings and document metadata with ChromaDB",
			ToolName:    "chromadb_collection_search",
			ToolDesc:    "Search ChromaDB embedded document collections",
			Triggers:    []string{"chromadb", "chroma", "local embeddings", "collection query"},
			Keywords:    []string{"chroma", "chromadb", "embeddings", "local", "rag"},
			Instructions: `# ChromaDB Embedded Vector Store
Manage lightweight local vector collections for RAG pipelines.`,
		},

		// === 8. Observability, Reliability & Incident Response ===
		{
			ID:          "prometheus_metrics_alert",
			Name:        "Prometheus Metrics & Alert Manager",
			Category:    "Observability",
			Description: "Queries Prometheus metrics, calculates P99 latency, and triggers alerts",
			ToolName:    "prometheus_query",
			ToolDesc:    "Query Prometheus time-series metrics",
			Triggers:    []string{"prometheus", "metrics", "p99", "latency", "alert", "promql"},
			Keywords:    []string{"prometheus", "promql", "metrics", "p99", "alerts"},
			Instructions: `# Prometheus Metrics & Alert Manager
Execute PromQL queries to measure system performance and service health.`,
		},
		{
			ID:          "datadog_apm_tracer",
			Name:        "Datadog APM & Trace Analyzer",
			Category:    "Observability",
			Description: "Inspects distributed flame graphs and pinpoints trace bottlenecks",
			ToolName:    "datadog_trace_inspect",
			ToolDesc:    "Inspect Datadog APM distributed traces",
			Triggers:    []string{"datadog", "apm", "trace", "flamegraph", "datadog trace"},
			Keywords:    []string{"datadog", "apm", "trace", "flamegraph", "spans"},
			Instructions: `# Datadog APM & Trace Analyzer
Inspect microservice call trees and detect slow database query spans.`,
		},
		{
			ID:          "sentry_error_tracker",
			Name:        "Sentry Real-Time Error Tracker",
			Category:    "Observability",
			Description: "Retrieves stack traces, unhandled exceptions, and breadcrumbs from Sentry",
			ToolName:    "sentry_get_issue",
			ToolDesc:    "Retrieve issue details and stack traces from Sentry",
			Triggers:    []string{"sentry", "stacktrace", "error", "exception", "sentry issue"},
			Keywords:    []string{"sentry", "error", "stacktrace", "exception", "breadcrumbs"},
			Instructions: `# Sentry Real-Time Error Tracker
Pinpoint code line numbers for production crashes and exceptions.`,
		},
		{
			ID:          "elastic_apm_collector",
			Name:        "Elastic APM Performance Monitor",
			Category:    "Observability",
			Description: "Tracks transaction spans and unhandled exceptions in Elastic APM",
			ToolName:    "elastic_apm_query",
			ToolDesc:    "Query Elastic APM spans and errors",
			Triggers:    []string{"elastic apm", "spans", "transactions", "elastic trace"},
			Keywords:    []string{"elastic", "apm", "spans", "transactions", "monitoring"},
			Instructions: `# Elastic APM Performance Monitor
Profile transaction latency across distributed application services.`,
		},
		{
			ID:          "pagerduty_incident_escalator",
			Name:        "PagerDuty On-Call Incident Escalator",
			Category:    "DevOps",
			Description: "Triggers on-call alerts, creates high-urgency incidents, and manages schedules",
			ToolName:    "pagerduty_trigger_incident",
			ToolDesc:    "Trigger PagerDuty on-call incident alerts",
			Triggers:    []string{"pagerduty", "oncall", "incident", "escalation", "pager", "trigger incident"},
			Keywords:    []string{"pagerduty", "incident", "oncall", "escalate", "alert"},
			Instructions: `# PagerDuty On-Call Incident Escalator
Escalate critical production outages to the active on-call engineer.`,
		},
		{
			ID:          "grafana_dashboard_manager",
			Name:        "Grafana Dashboard & Alert Manager",
			Category:    "Observability",
			Description: "Provisions Grafana dashboard JSON models, configures panel alert thresholds, and links data sources",
			ToolName:    "grafana_manage_dashboard",
			ToolDesc:    "Provision and update Grafana dashboards and alerts",
			Triggers:    []string{"grafana", "dashboard json", "grafana panel", "alert threshold"},
			Keywords:    []string{"grafana", "dashboard", "panel", "visualization", "alerts"},
			Instructions: `# Grafana Dashboard & Alert Manager
Maintain infrastructure telemetry visual panels and alert rules.`,
		},

		// === 9. API, Microservices & Communications ===
		{
			ID:          "graphql_api_client",
			Name:        "GraphQL Schema & Query Executor",
			Category:    "API",
			Description: "Introspects GraphQL schemas and executes GraphQL mutations and queries",
			ToolName:    "graphql_execute",
			ToolDesc:    "Execute GraphQL queries and mutations",
			Triggers:    []string{"graphql", "mutation", "graphql query", "graphql schema"},
			Keywords:    []string{"graphql", "query", "mutation", "schema", "api"},
			Instructions: `# GraphQL Schema & Query Executor
Execute strongly-typed GraphQL operations against backend services.`,
		},
		{
			ID:          "graphql_apollo_federation",
			Name:        "Apollo Federation Subgraph Router",
			Category:    "API",
			Description: "Composes subgraphs into a unified supergraph gateway schema",
			ToolName:    "apollo_compose_subgraph",
			ToolDesc:    "Validate and compose Apollo Federation subgraphs",
			Triggers:    []string{"apollo", "federation", "supergraph", "subgraph", "apollo router"},
			Keywords:    []string{"apollo", "federation", "supergraph", "subgraph", "graphql"},
			Instructions: `# Apollo Federation Subgraph Router
Validate entity resolvers and compose distributed microservice graphs.`,
		},
		{
			ID:          "grpc_protobuf_caller",
			Name:        "gRPC Protobuf Service Caller",
			Category:    "API",
			Description: "Encodes protobuf payloads and makes high-speed gRPC RPC calls",
			ToolName:    "grpc_rpc_invoke",
			ToolDesc:    "Invoke gRPC service methods with protobuf",
			Triggers:    []string{"grpc", "protobuf", "rpc", "proto", "grpc call"},
			Keywords:    []string{"grpc", "protobuf", "rpc", "proto", "service"},
			Instructions: `# gRPC Protobuf Service Caller
Invoke binary gRPC endpoints with compiled protocol buffers.`,
		},
		{
			ID:          "fastapi_openapi_tester",
			Name:        "FastAPI OpenAPI Endpoint Tester",
			Category:    "API",
			Description: "Tests REST endpoints against auto-generated Swagger OpenAPI specs",
			ToolName:    "fastapi_test_route",
			ToolDesc:    "Test FastAPI routes and OpenAPI schemas",
			Triggers:    []string{"fastapi", "openapi", "swagger", "rest api", "endpoint test"},
			Keywords:    []string{"fastapi", "openapi", "swagger", "rest", "api"},
			Instructions: `# FastAPI OpenAPI Endpoint Tester
Validate HTTP REST schemas and execute contract tests.`,
		},
		{
			ID:          "slack_webhook_notifier",
			Name:        "Slack Webhook & Alert Bot",
			Category:    "Communication",
			Description: "Dispatches rich Block Kit notifications and incident alerts to Slack",
			ToolName:    "slack_send_message",
			ToolDesc:    "Send rich notification messages to Slack channels",
			Triggers:    []string{"slack", "notification", "alert", "channel", "message", "slack webhook"},
			Keywords:    []string{"slack", "message", "webhook", "block kit", "channel"},
			Instructions: `# Slack Webhook & Alert Bot
Broadcast interactive Block Kit cards to team collaboration channels.`,
		},
		{
			ID:          "sendgrid_email_sender",
			Name:        "SendGrid Transactional Email Sender",
			Category:    "Communication",
			Description: "Dispatches password resets and transactional email receipts via SendGrid",
			ToolName:    "sendgrid_send_email",
			ToolDesc:    "Send transactional emails via SendGrid API",
			Triggers:    []string{"sendgrid", "email", "transactional email", "receipt email", "send email"},
			Keywords:    []string{"sendgrid", "email", "smtp", "transactional", "mail"},
			Instructions: `# SendGrid Transactional Email Sender
Deliver templated notifications and verification emails at scale.`,
		},
		{
			ID:          "twillio_sms_notifier",
			Name:        "Twilio SMS & 2FA Auth Dispatcher",
			Category:    "Communication",
			Description: "Dispatches SMS verification codes and 2FA alerts via Twilio API",
			ToolName:    "twilio_send_sms",
			ToolDesc:    "Send SMS text messages via Twilio API",
			Triggers:    []string{"twilio", "sms", "2fa", "text message", "send sms"},
			Keywords:    []string{"twilio", "sms", "2fa", "text", "phone"},
			Instructions: `# Twilio SMS & 2FA Auth Dispatcher
Send two-factor authentication codes and priority SMS notices.`,
		},
		{
			ID:          "discord_bot_dispatcher",
			Name:        "Discord Webhook & Bot Dispatcher",
			Category:    "Communication",
			Description: "Dispatches rich Discord embed cards, monitors server channels, and executes bot interactions",
			ToolName:    "discord_dispatch_embed",
			ToolDesc:    "Send rich embed messages to Discord channels",
			Triggers:    []string{"discord", "discord webhook", "discord embed", "discord bot"},
			Keywords:    []string{"discord", "bot", "embed", "webhook", "community"},
			Instructions: `# Discord Webhook & Bot Dispatcher
Format and publish rich community updates and build status embeds.`,
		},

		// === 10. Security, Secrets & Feature Management ===
		{
			ID:          "vault_secret_manager",
			Name:        "HashiCorp Vault Secret Manager",
			Category:    "Security",
			Description: "Reads and writes dynamic secrets and encryption keys from Vault",
			ToolName:    "vault_secret_read",
			ToolDesc:    "Read and write HashiCorp Vault secrets",
			Triggers:    []string{"vault", "secret", "hashicorp", "encryption", "vault secret"},
			Keywords:    []string{"vault", "secret", "hashicorp", "token", "encryption"},
			Instructions: `# HashiCorp Vault Secret Manager
Fetch and rotate sensitive credentials from HashiCorp Vault KV engines.`,
		},
		{
			ID:          "iam_security_guard",
			Name:        "IAM Security & Token Rotator",
			Category:    "Security",
			Description: "Rotates OAuth2 service account tokens, audits access logs, and validates HMAC signatures",
			ToolName:    "iam_auth_rotator",
			ToolDesc:    "Rotate IAM credentials and OAuth tokens",
			Triggers:    []string{"iam", "oauth2", "token", "rotate", "security", "credentials"},
			Keywords:    []string{"iam", "oauth2", "token", "security", "auth"},
			Instructions: `# IAM Security & Token Rotator
Audit permissions and automate access token rotation cycles.`,
		},
		{
			ID:          "launchdarkly_feature_flags",
			Name:        "LaunchDarkly Feature Flag Controller",
			Category:    "DevOps",
			Description: "Toggles feature flags and manages rollout percentage rules",
			ToolName:    "launchdarkly_toggle_flag",
			ToolDesc:    "Toggle LaunchDarkly feature flags",
			Triggers:    []string{"launchdarkly", "feature flag", "toggle", "rollout", "kill switch"},
			Keywords:    []string{"launchdarkly", "feature", "flag", "rollout", "toggle"},
			Instructions: `# LaunchDarkly Feature Flag Controller
Manage progressive rollouts and dynamic application feature gates.`,
		},
		{
			ID:          "snyk_vulnerability_scanner",
			Name:        "Snyk Dependency & Container Scanner",
			Category:    "Security",
			Description: "Scans open source dependencies, npm/pypi lockfiles, and container base images for known CVE vulnerabilities",
			ToolName:    "snyk_scan_dependencies",
			ToolDesc:    "Scan software dependencies for CVE security vulnerabilities with Snyk",
			Triggers:    []string{"snyk", "cve scan", "dependency vulnerability", "npm audit", "container vulnerability"},
			Keywords:    []string{"snyk", "cve", "vulnerability", "security", "dependencies"},
			Instructions: `# Snyk Dependency & Container Scanner
Identify and remediate vulnerable third-party libraries and container packages.`,
		},
		{
			ID:          "linear_project_tracker",
			Name:        "Linear Issue & Project Tracker",
			Category:    "Management",
			Description: "Creates Linear issue tickets, assigns cycle sprints, updates priority states, and manages project milestones",
			ToolName:    "linear_manage_issue",
			ToolDesc:    "Create and update Linear project issues and cycle sprints",
			Triggers:    []string{"linear", "linear issue", "sprint cycle", "project roadmap", "ticket status"},
			Keywords:    []string{"linear", "issue", "sprint", "project", "cycle", "tracker"},
			Instructions: `# Linear Issue & Project Tracker
Track software development milestones and sprint backlog tickets in Linear.`,
		},
		{
			ID:          "jira_issue_tracker",
			Name:        "Jira Issue & Sprint Board Manager",
			Category:    "Management",
			Description: "Creates tickets, updates sprint story points, and assigns bugs in Jira",
			ToolName:    "jira_update_ticket",
			ToolDesc:    "Update Jira tickets and sprint boards",
			Triggers:    []string{"jira", "ticket", "sprint", "bug", "issue", "jira ticket"},
			Keywords:    []string{"jira", "ticket", "sprint", "bug", "issue"},
			Instructions: `# Jira Issue & Sprint Board Manager
Update agile sprint backlogs, bug reports, and workflow transitions in Jira.`,
		},
		{
			ID:          "pytest_test_runner",
			Name:        "PyTest Automated Test Suite Runner",
			Category:    "Testing",
			Description: "Executes unit and integration test fixtures and parses coverage reports",
			ToolName:    "pytest_run_suite",
			ToolDesc:    "Run PyTest test suite with coverage",
			Triggers:    []string{"pytest", "test runner", "coverage", "fixtures", "python test"},
			Keywords:    []string{"pytest", "test", "coverage", "python", "unit"},
			Instructions: `# PyTest Automated Test Suite Runner
Execute Python test fixtures and generate code coverage reports.`,
		},
	}
}

// PopulateCatalog registers the eval catalog in memory.
// If targetDir is non-empty it also writes markdown copies there; tests should pass "".
func PopulateCatalog(reg *SkillRegistry, targetDir string) error {
	if targetDir != "" {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target skills dir: %w", err)
		}
	}

	catalog := GetUniversalSkillCatalog()
	for _, item := range catalog {
		toolProps := item.Parameters.Properties
		if toolProps == nil {
			toolProps = map[string]ParameterProperty{
				"query": {
					Type:        "string",
					Description: fmt.Sprintf("Primary query or action for %s", item.ToolName),
				},
				"options": {
					Type:        "string",
					Description: "Optional execution parameters (JSON string)",
				},
			}
		}

		skill := &Skill{
			ID:          item.ID,
			Name:        item.Name,
			Category:    item.Category,
			Description: item.Description,
			Triggers:    item.Triggers,
			Keywords:    item.Keywords,
			Tools: []ToolDefinition{
				{
					Name:        item.ToolName,
					Description: item.ToolDesc,
					Parameters: ToolParameters{
						Type:       "object",
						Properties: toolProps,
						Required:   []string{"query"},
					},
					Category: item.Category,
				},
			},
			Instructions: item.Instructions,
			CreatedAt:    time.Now(),
		}

		reg.Register(skill)

		// Persist as Markdown file if targetDir provided
		if targetDir != "" {
			filePath := filepath.Join(targetDir, item.ID+".md")
			triggersJSON, _ := json.Marshal(item.Triggers)
			keywordsJSON, _ := json.Marshal(item.Keywords)

			mdContent := fmt.Sprintf(`---
name: "%s"
id: "%s"
category: "%s"
execution: "sandbox-mock"
description: "%s"
triggers: %s
keywords: %s
---

# %s

> Eval fixture. Sandbox mock — not a live vendor API.

%s

`+"```json"+`
{
  "name": "%s",
  "description": "%s",
  "parameters": {
    "type": "object",
    "properties": {
      "query": {"type": "string", "description": "Primary action query"},
      "options": {"type": "string", "description": "Optional parameters"}
    },
    "required": ["query"]
  }
}
`+"```"+`
`, item.Name, item.ID, item.Category, item.Description, string(triggersJSON), string(keywordsJSON), item.Name, item.Description, item.ToolName, item.ToolDesc)

			_ = os.WriteFile(filePath, []byte(mdContent), 0644)
		}
	}

	return nil
}
