// Package ilicense provides offline license activation and validation for Go applications.
//
// The typical flow is:
//  1. Build a Client with DefaultConfig or custom Config.
//  2. Call Init() at application startup.
//  3. Activate once with an activation code issued by your management platform.
//  4. Guard protected features by calling CheckLicense or CheckModule.
//
// This package intentionally exposes only stable SDK APIs. Internal verification
// details are kept in non-public packages.
package ilicense
