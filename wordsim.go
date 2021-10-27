package main

import (
	"context"
	"errors"
	"github.com/zone-7/andflow_plugin"
	"strconv"
	"strings"
)
type Andflow_plugin_wordsim struct {

}
func (a *Andflow_plugin_wordsim)GetName() string{
	return "wordsim"
}
func (a *Andflow_plugin_wordsim)Init(callback interface{}){

}
func (a *Andflow_plugin_wordsim)PrepareMetadata(userid int,flowCode string, metadata string)string{
	md:=andflow_plugin.ParseMetadata(metadata)
	md.Name=a.GetName()
	md.Title="文本相似度"
	md.Group="自然语言"
	md.Tag="机器学习"
	md.Params=[]andflow_plugin.MetadataPropertiesModel{
		andflow_plugin.MetadataPropertiesModel{Name:"text1",Title:"文本1参数",Default:"text1"},
		andflow_plugin.MetadataPropertiesModel{Name:"text2",Title:"文本2参数",Default:"text2"},
		andflow_plugin.MetadataPropertiesModel{Name:"method",Title:"方法",Default:"cos",
			Element:"select" ,Options: []andflow_plugin.MetadataOptionModel{
			andflow_plugin.MetadataOptionModel{Value:"cos",Label:"余弦"},
			andflow_plugin.MetadataOptionModel{Value:"hash",Label:"哈希"},

			}},

		andflow_plugin.MetadataPropertiesModel{Name:"tops",Title:"高频词个数",Placeholder:"多少个头部高频词用于比较",Default:"30"},

		andflow_plugin.MetadataPropertiesModel{Name:"dict",Title:"词语列表",Placeholder:"逗号分隔多个词语"},
		andflow_plugin.MetadataPropertiesModel{Name:"result",Title:"比对结果参数",Default:"result"},

	}
	return md.ToJson()

}


func (a *Andflow_plugin_wordsim)Filter(ctx context.Context,runtimeId string,preActionId string, actionId string,callback interface{})(bool,error){

	return true,nil
}


func (a *Andflow_plugin_wordsim)Exec(ctx context.Context,runtimeId string,preActionId string, actionId string,callback interface{})(interface{},error){

	actionCallback := andflow_plugin.ParseActionCallbacker(callback)
	key_text1 := actionCallback.GetActionParam(actionId, "text1")
	key_text2 := actionCallback.GetActionParam(actionId, "text2")
	dict := actionCallback.GetActionParam(actionId, "dict")

	tops := actionCallback.GetActionParam(actionId, "tops")

	method := actionCallback.GetActionParam(actionId, "method")

	key_result := actionCallback.GetActionParam(actionId, "result")

	text1:=actionCallback.GetRuntimeParam(key_text1)
	text2:=actionCallback.GetRuntimeParam(key_text2)
	srcStr,ok1:=text1.(string)
	dstStr,ok2:=text2.(string)
	if !ok1 || !ok2{
		return nil,errors.New("输入的文本不是字符串")
	}


	g := NewGoJieba()

	if len(dict)>0{
		commonWords:=make([]string,0)
		words := strings.Split(dict,",")
		for _,word:=range words{
			if len(word)==0{
				continue
			}
			w:=strings.Split(word,"，")
			if len(w)==0{
				continue
			}
			commonWords = append(commonWords,w...)
		}

		g.AddWords(commonWords)
	}

 	srcStr = removeHtml(srcStr)
	dstStr = removeHtml(dstStr)

	if method=="cos"{
		srcWords := g.C.Cut(srcStr, true)
		dstWords := g.C.Cut(dstStr, true)
		score := CosineSimilar(srcWords, dstWords)
		actionCallback.SetRuntimeParam(key_result, score)

	} else  {

		if len(tops)==0{
			tops="30"
		}
		topk := 30
		topCount,err := strconv.Atoi(tops)
		if err==nil{
			topk= topCount
		}


		srcWordsWeight := g.C.ExtractWithWeight(srcStr, topk)
		dstWordsWeight := g.C.ExtractWithWeight(dstStr, topk)

		distance, err := SimHashSimilar(srcWordsWeight, dstWordsWeight)
		actionCallback.SetRuntimeParam(key_result, distance)

	}







	return nil,nil
}


func main(){
	andflow_plugin.InitPlugin(&Andflow_plugin_wordsim{})
}