package dotgit

import (
	"os"
)

func (s *SuiteDotGit) TestRepositoryFilesystem() {
	fs := s.EmptyFS()

	err := fs.MkdirAll("dotGit", 0777)
	s.Require().NoError(err)
	dotGitFs, err := fs.Chroot("dotGit")
	s.Require().NoError(err)

	err = fs.MkdirAll("commonDotGit", 0777)
	s.Require().NoError(err)
	commonDotGitFs, err := fs.Chroot("commonDotGit")
	s.Require().NoError(err)

	repositoryFs := NewRepositoryFilesystem(dotGitFs, commonDotGitFs)
	s.Equal(dotGitFs.Root(), repositoryFs.Root())

	somedir, err := repositoryFs.Chroot("somedir")
	s.Require().NoError(err)
	s.Equal(repositoryFs.Join(dotGitFs.Root(), "somedir"), somedir.Root())

	_, err = repositoryFs.Create("somefile")
	s.Require().NoError(err)

	_, err = repositoryFs.Stat("somefile")
	s.Require().NoError(err)

	file, err := repositoryFs.Open("somefile")
	s.Require().NoError(err)
	err = file.Close()
	s.Require().NoError(err)

	file, err = repositoryFs.OpenFile("somefile", os.O_RDONLY, 0666)
	s.Require().NoError(err)
	err = file.Close()
	s.Require().NoError(err)

	file, err = repositoryFs.Create("somefile2")
	s.Require().NoError(err)
	err = file.Close()
	s.Require().NoError(err)
	_, err = repositoryFs.Stat("somefile2")
	s.Require().NoError(err)
	err = repositoryFs.Rename("somefile2", "newfile")
	s.Require().NoError(err)

	tempDir, err := repositoryFs.TempFile("tmp", "myprefix")
	s.Require().NoError(err)
	s.Equal(repositoryFs.Join(dotGitFs.Root(), "tmp", tempDir.Name()), repositoryFs.Join(repositoryFs.Root(), "tmp", tempDir.Name()))

	err = repositoryFs.Symlink("newfile", "somelink")
	s.Require().NoError(err)

	_, err = repositoryFs.Lstat("somelink")
	s.Require().NoError(err)

	link, err := repositoryFs.Readlink("somelink")
	s.Require().NoError(err)
	s.Equal("newfile", link)

	err = repositoryFs.Remove("somelink")
	s.Require().NoError(err)

	_, err = repositoryFs.Stat("somelink")
	s.True(os.IsNotExist(err))

	dirs := []string{objectsPath, refsPath, packedRefsPath, configPath, branchesPath, hooksPath, infoPath, remotesPath, logsPath, shallowPath, worktreesPath}
	for _, dir := range dirs {
		err := repositoryFs.MkdirAll(dir, 0777)
		s.Require().NoError(err)
		_, err = commonDotGitFs.Stat(dir)
		s.Require().NoError(err)
		_, err = dotGitFs.Stat(dir)
		s.True(os.IsNotExist(err))
	}

	exceptionsPaths := []string{repositoryFs.Join(logsPath, "HEAD"), repositoryFs.Join(refsPath, "bisect"), repositoryFs.Join(refsPath, "rewritten"), repositoryFs.Join(refsPath, "worktree")}
	for _, path := range exceptionsPaths {
		_, err := repositoryFs.Create(path)
		s.Require().NoError(err)
		_, err = commonDotGitFs.Stat(path)
		s.True(os.IsNotExist(err))
		_, err = dotGitFs.Stat(path)
		s.Require().NoError(err)
	}

	err = repositoryFs.MkdirAll("refs/heads", 0777)
	s.Require().NoError(err)
	_, err = commonDotGitFs.Stat("refs/heads")
	s.Require().NoError(err)
	_, err = dotGitFs.Stat("refs/heads")
	s.True(os.IsNotExist(err))

	err = repositoryFs.MkdirAll("objects/pack", 0777)
	s.Require().NoError(err)
	_, err = commonDotGitFs.Stat("objects/pack")
	s.Require().NoError(err)
	_, err = dotGitFs.Stat("objects/pack")
	s.True(os.IsNotExist(err))

	err = repositoryFs.MkdirAll("a/b/c", 0777)
	s.Require().NoError(err)
	_, err = commonDotGitFs.Stat("a/b/c")
	s.True(os.IsNotExist(err))
	_, err = dotGitFs.Stat("a/b/c")
	s.Require().NoError(err)
}
