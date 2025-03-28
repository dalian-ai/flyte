package update

import (
	"fmt"
	"testing"

	"github.com/flyteorg/flyte/flytectl/cmd/config/subcommand/workflowexecutionconfig"
	"github.com/flyteorg/flyte/flytectl/cmd/testutils"
	"github.com/flyteorg/flyte/flytectl/pkg/ext"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	validProjectWorkflowExecutionConfigFilePath       = "testdata/valid_project_workflow_execution_config.yaml"
	validProjectDomainWorkflowExecutionConfigFilePath = "testdata/valid_project_domain_workflow_execution_config.yaml"
	validWorkflowExecutionConfigFilePath              = "testdata/valid_workflow_workflow_execution_config.yaml"
)

func TestWorkflowExecutionConfigUpdateRequiresAttributeFile(t *testing.T) {
	testWorkflowExecutionConfigUpdate(
		t,
		/* setup */ nil,
		/* assert */ func(s *testutils.TestStruct, err error) {
			assert.ErrorContains(t, err, "attrFile is mandatory")
			s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
}

func TestWorkflowExecutionConfigUpdateFailsWhenAttributeFileDoesNotExist(t *testing.T) {
	testWorkflowExecutionConfigUpdate(
		t,
		/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
			config.AttrFile = testDataNonExistentFile
			config.Force = true
		},
		/* assert */ func(s *testutils.TestStruct, err error) {
			assert.ErrorContains(t, err, "unable to read from testdata/non-existent-file yaml file")
			s.UpdaterExt.AssertNotCalled(t, "FetchWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
}

func TestWorkflowExecutionConfigUpdateFailsWhenAttributeFileIsMalformed(t *testing.T) {
	testWorkflowExecutionConfigUpdate(
		t,
		/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
			config.AttrFile = testDataInvalidAttrFile
			config.Force = true
		},
		/* assert */ func(s *testutils.TestStruct, err error) {
			assert.ErrorContains(t, err, "error unmarshaling JSON: while decoding JSON: json: unknown field \"InvalidDomain\"")
			s.UpdaterExt.AssertNotCalled(t, "FetchWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
}

func TestWorkflowExecutionConfigUpdateHappyPath(t *testing.T) {
	t.Run("workflow", func(t *testing.T) {
		testWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				s.TearDownAndVerifyContains(t, `Updated attributes from flytesnacks project and domain development`)
			})
	})

	t.Run("domain", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				s.TearDownAndVerifyContains(t, `Updated attributes from flytesnacks project and domain development`)
			})
	})

	t.Run("project", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
				s.TearDownAndVerifyContains(t, `Updated attributes from flytesnacks project`)
			})
	})
}

func TestWorkflowExecutionConfigUpdateFailsWithoutForceFlag(t *testing.T) {
	t.Run("workflow", func(t *testing.T) {
		testWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.Force = false
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.ErrorContains(t, err, "update aborted by user")
				s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("domain", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.Force = false
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.ErrorContains(t, err, "update aborted by user")
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("project", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.Force = false
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.ErrorContains(t, err, "update aborted by user")
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
			})
	})
}

func TestWorkflowExecutionConfigUpdateDoesNothingWithDryRunFlag(t *testing.T) {
	t.Run("workflow", func(t *testing.T) {
		testWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("domain", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("project", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
			})
	})
}

func TestWorkflowExecutionConfigUpdateIgnoresForceFlagWithDryRun(t *testing.T) {
	t.Run("workflow without --force", func(t *testing.T) {
		testWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.Force = false
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("workflow with --force", func(t *testing.T) {
		testWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.Force = true
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("domain without --force", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.Force = false
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("domain with --force", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.Force = true
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("project without --force", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.Force = false
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("project with --force", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdate(
			t,
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.Force = true
				config.DryRun = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertNotCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
			})
	})
}

func TestWorkflowExecutionConfigUpdateSucceedsWhenAttributesDoNotExist(t *testing.T) {
	t.Run("workflow", func(t *testing.T) {
		testWorkflowExecutionConfigUpdateWithMockSetup(
			t,
			/* mockSetup */ func(s *testutils.TestStruct, target *admin.WorkflowAttributes) {
				s.FetcherExt.
					EXPECT().FetchWorkflowAttributes(s.Ctx, target.GetProject(), target.GetDomain(), target.GetWorkflow(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
					Return(nil, ext.NewNotFoundError("attribute"))
				s.UpdaterExt.
					EXPECT().UpdateWorkflowAttributes(s.Ctx, target.GetProject(), target.GetDomain(), target.GetWorkflow(), mock.Anything).
					Return(nil)
			},
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				s.TearDownAndVerifyContains(t, `Updated attributes from flytesnacks project and domain development`)
			})
	})

	t.Run("domain", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdateWithMockSetup(
			t,
			/* mockSetup */ func(s *testutils.TestStruct, target *admin.ProjectDomainAttributes) {
				s.FetcherExt.
					EXPECT().FetchProjectDomainAttributes(s.Ctx, target.GetProject(), target.GetDomain(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
					Return(nil, ext.NewNotFoundError("attribute"))
				s.UpdaterExt.
					EXPECT().UpdateProjectDomainAttributes(s.Ctx, target.GetProject(), target.GetDomain(), mock.Anything).
					Return(nil)
			},
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				s.TearDownAndVerifyContains(t, `Updated attributes from flytesnacks project and domain development`)
			})
	})

	t.Run("project", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdateWithMockSetup(
			t,
			/* mockSetup */ func(s *testutils.TestStruct, target *admin.ProjectAttributes) {
				s.FetcherExt.
					EXPECT().FetchProjectAttributes(s.Ctx, target.GetProject(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
					Return(nil, ext.NewNotFoundError("attribute"))
				s.UpdaterExt.
					EXPECT().UpdateProjectAttributes(s.Ctx, target.GetProject(), mock.Anything).
					Return(nil)
			},
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Nil(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
				s.TearDownAndVerifyContains(t, `Updated attributes from flytesnacks project`)
			})
	})
}

func TestWorkflowExecutionConfigUpdateFailsWhenAdminClientFails(t *testing.T) {
	t.Run("workflow", func(t *testing.T) {
		testWorkflowExecutionConfigUpdateWithMockSetup(
			t,
			/* mockSetup */ func(s *testutils.TestStruct, target *admin.WorkflowAttributes) {
				s.FetcherExt.
					EXPECT().FetchWorkflowAttributes(s.Ctx, target.GetProject(), target.GetDomain(), target.GetWorkflow(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
					Return(&admin.WorkflowAttributesGetResponse{Attributes: target}, nil)
				s.UpdaterExt.
					EXPECT().UpdateWorkflowAttributes(s.Ctx, target.GetProject(), target.GetDomain(), target.GetWorkflow(), mock.Anything).
					Return(fmt.Errorf("network error"))
			},
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes) {
				config.AttrFile = validWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Error(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateWorkflowAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("domain", func(t *testing.T) {
		testProjectDomainWorkflowExecutionConfigUpdateWithMockSetup(
			t,
			/* mockSetup */ func(s *testutils.TestStruct, target *admin.ProjectDomainAttributes) {
				s.FetcherExt.
					EXPECT().FetchProjectDomainAttributes(s.Ctx, target.GetProject(), target.GetDomain(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
					Return(&admin.ProjectDomainAttributesGetResponse{Attributes: target}, nil)
				s.UpdaterExt.
					EXPECT().UpdateProjectDomainAttributes(s.Ctx, target.GetProject(), target.GetDomain(), mock.Anything).
					Return(fmt.Errorf("network error"))
			},
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes) {
				config.AttrFile = validProjectDomainWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Error(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateProjectDomainAttributes", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
	})

	t.Run("project", func(t *testing.T) {
		testProjectWorkflowExecutionConfigUpdateWithMockSetup(
			t,
			/* mockSetup */ func(s *testutils.TestStruct, target *admin.ProjectAttributes) {
				s.FetcherExt.
					EXPECT().FetchProjectAttributes(s.Ctx, target.GetProject(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
					Return(&admin.ProjectAttributesGetResponse{Attributes: target}, nil)
				s.UpdaterExt.
					EXPECT().UpdateProjectAttributes(s.Ctx, target.GetProject(), mock.Anything).
					Return(fmt.Errorf("network error"))
			},
			/* setup */ func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes) {
				config.AttrFile = validProjectWorkflowExecutionConfigFilePath
				config.Force = true
			},
			/* assert */ func(s *testutils.TestStruct, err error) {
				assert.Error(t, err)
				s.UpdaterExt.AssertCalled(t, "UpdateProjectAttributes", mock.Anything, mock.Anything, mock.Anything)
			})
	})
}

func testWorkflowExecutionConfigUpdate(
	t *testing.T,
	setup func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes),
	asserter func(s *testutils.TestStruct, err error),
) {
	testWorkflowExecutionConfigUpdateWithMockSetup(
		t,
		/* mockSetup */ func(s *testutils.TestStruct, target *admin.WorkflowAttributes) {
			s.FetcherExt.
				EXPECT().FetchWorkflowAttributes(s.Ctx, target.GetProject(), target.GetDomain(), target.GetWorkflow(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
				Return(&admin.WorkflowAttributesGetResponse{Attributes: target}, nil)
			s.UpdaterExt.
				EXPECT().UpdateWorkflowAttributes(s.Ctx, target.GetProject(), target.GetDomain(), target.GetWorkflow(), mock.Anything).
				Return(nil)
		},
		setup,
		asserter,
	)
}

func testWorkflowExecutionConfigUpdateWithMockSetup(
	t *testing.T,
	mockSetup func(s *testutils.TestStruct, target *admin.WorkflowAttributes),
	setup func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.WorkflowAttributes),
	asserter func(s *testutils.TestStruct, err error),
) {
	s := testutils.Setup(t)

	workflowexecutionconfig.DefaultUpdateConfig = &workflowexecutionconfig.AttrUpdateConfig{}
	target := newTestWorkflowExecutionConfig()

	if mockSetup != nil {
		mockSetup(&s, target)
	}

	if setup != nil {
		setup(&s, workflowexecutionconfig.DefaultUpdateConfig, target)
	}

	err := updateWorkflowExecutionConfigFunc(s.Ctx, nil, s.CmdCtx)

	if asserter != nil {
		asserter(&s, err)
	}

	// cleanup
	workflowexecutionconfig.DefaultUpdateConfig = &workflowexecutionconfig.AttrUpdateConfig{}
}

func newTestWorkflowExecutionConfig() *admin.WorkflowAttributes {
	return &admin.WorkflowAttributes{
		// project, domain, and workflow names need to be same as in the tests spec files in testdata folder
		Project:  "flytesnacks",
		Domain:   "development",
		Workflow: "core.control_flow.merge_sort.merge_sort",
		MatchingAttributes: &admin.MatchingAttributes{
			Target: &admin.MatchingAttributes_WorkflowExecutionConfig{
				WorkflowExecutionConfig: &admin.WorkflowExecutionConfig{
					MaxParallelism: 1337,
					Annotations: &admin.Annotations{
						Values: map[string]string{
							testutils.RandomName(5): testutils.RandomName(10),
							testutils.RandomName(5): testutils.RandomName(10),
							testutils.RandomName(5): testutils.RandomName(10),
						},
					},
				},
			},
		},
	}
}

func testProjectWorkflowExecutionConfigUpdate(
	t *testing.T,
	setup func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes),
	asserter func(s *testutils.TestStruct, err error),
) {
	testProjectWorkflowExecutionConfigUpdateWithMockSetup(
		t,
		/* mockSetup */ func(s *testutils.TestStruct, target *admin.ProjectAttributes) {
			s.FetcherExt.
				EXPECT().FetchProjectAttributes(s.Ctx, target.GetProject(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
				Return(&admin.ProjectAttributesGetResponse{Attributes: target}, nil)
			s.UpdaterExt.
				EXPECT().UpdateProjectAttributes(s.Ctx, target.GetProject(), mock.Anything).
				Return(nil)
		},
		setup,
		asserter,
	)
}

func testProjectWorkflowExecutionConfigUpdateWithMockSetup(
	t *testing.T,
	mockSetup func(s *testutils.TestStruct, target *admin.ProjectAttributes),
	setup func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectAttributes),
	asserter func(s *testutils.TestStruct, err error),
) {
	s := testutils.Setup(t)

	workflowexecutionconfig.DefaultUpdateConfig = &workflowexecutionconfig.AttrUpdateConfig{}
	target := newTestProjectWorkflowExecutionConfig()

	if mockSetup != nil {
		mockSetup(&s, target)
	}

	if setup != nil {
		setup(&s, workflowexecutionconfig.DefaultUpdateConfig, target)
	}

	err := updateWorkflowExecutionConfigFunc(s.Ctx, nil, s.CmdCtx)

	if asserter != nil {
		asserter(&s, err)
	}

	// cleanup
	workflowexecutionconfig.DefaultUpdateConfig = &workflowexecutionconfig.AttrUpdateConfig{}
}

func newTestProjectWorkflowExecutionConfig() *admin.ProjectAttributes {
	return &admin.ProjectAttributes{
		// project name needs to be same as in the tests spec files in testdata folder
		Project: "flytesnacks",
		MatchingAttributes: &admin.MatchingAttributes{
			Target: &admin.MatchingAttributes_WorkflowExecutionConfig{
				WorkflowExecutionConfig: &admin.WorkflowExecutionConfig{
					MaxParallelism: 1337,
					Annotations: &admin.Annotations{
						Values: map[string]string{
							testutils.RandomName(5): testutils.RandomName(10),
							testutils.RandomName(5): testutils.RandomName(10),
							testutils.RandomName(5): testutils.RandomName(10),
						},
					},
				},
			},
		},
	}
}

func testProjectDomainWorkflowExecutionConfigUpdate(
	t *testing.T,
	setup func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes),
	asserter func(s *testutils.TestStruct, err error),
) {
	testProjectDomainWorkflowExecutionConfigUpdateWithMockSetup(
		t,
		/* mockSetup */ func(s *testutils.TestStruct, target *admin.ProjectDomainAttributes) {
			s.FetcherExt.
				EXPECT().FetchProjectDomainAttributes(s.Ctx, target.GetProject(), target.GetDomain(), admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG).
				Return(&admin.ProjectDomainAttributesGetResponse{Attributes: target}, nil)
			s.UpdaterExt.
				EXPECT().UpdateProjectDomainAttributes(s.Ctx, target.GetProject(), target.GetDomain(), mock.Anything).
				Return(nil)
		},
		setup,
		asserter,
	)
}

func testProjectDomainWorkflowExecutionConfigUpdateWithMockSetup(
	t *testing.T,
	mockSetup func(s *testutils.TestStruct, target *admin.ProjectDomainAttributes),
	setup func(s *testutils.TestStruct, config *workflowexecutionconfig.AttrUpdateConfig, target *admin.ProjectDomainAttributes),
	asserter func(s *testutils.TestStruct, err error),
) {
	s := testutils.Setup(t)

	workflowexecutionconfig.DefaultUpdateConfig = &workflowexecutionconfig.AttrUpdateConfig{}
	target := newTestProjectDomainWorkflowExecutionConfig()

	if mockSetup != nil {
		mockSetup(&s, target)
	}

	if setup != nil {
		setup(&s, workflowexecutionconfig.DefaultUpdateConfig, target)
	}

	err := updateWorkflowExecutionConfigFunc(s.Ctx, nil, s.CmdCtx)

	if asserter != nil {
		asserter(&s, err)
	}

	// cleanup
	workflowexecutionconfig.DefaultUpdateConfig = &workflowexecutionconfig.AttrUpdateConfig{}
}

func newTestProjectDomainWorkflowExecutionConfig() *admin.ProjectDomainAttributes {
	return &admin.ProjectDomainAttributes{
		// project and domain names need to be same as in the tests spec files in testdata folder
		Project: "flytesnacks",
		Domain:  "development",
		MatchingAttributes: &admin.MatchingAttributes{
			Target: &admin.MatchingAttributes_WorkflowExecutionConfig{
				WorkflowExecutionConfig: &admin.WorkflowExecutionConfig{
					MaxParallelism: 1337,
					Annotations: &admin.Annotations{
						Values: map[string]string{
							testutils.RandomName(5): testutils.RandomName(10),
							testutils.RandomName(5): testutils.RandomName(10),
							testutils.RandomName(5): testutils.RandomName(10),
						},
					},
				},
			},
		},
	}
}
