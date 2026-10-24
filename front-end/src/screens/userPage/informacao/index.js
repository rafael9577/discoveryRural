import React, { useState, useEffect } from "react";
import { ImageBackground } from "react-native";
import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, FlatList} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import style from "./style";
import { listarAnimais } from "../../../../services/API"; 

export default () => {
  const [animais, setAnimais] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const carregarAnimais = async () => {
    try {
      setLoading(true);
      setError(null);
      const dados = await listarAnimais();
      console.log(dados);
      setAnimais(dados);
    } catch (err) {
      console.error("Erro ao carregar os animais:", err);
      setError("Não foi possível carregar os animais.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    carregarAnimais();
  }, []);

  const renderCardItem = ({ item }) => (
    <View style={style.cardContainer}>
      <Image
        source={{ 
          uri: item.imagem_url || "https://via.placeholder.com/300x180?text=Sem+Imagem" 
        }}
        style={style.cardImage}
        resizeMode="cover"
      />
      <View>
        <TouchableOpacity style={""} activeOpacity={0.8}>
          <View style={{flex: 1, alignContent:"flex-start"}}>
              <Text style={style.cardTitle}>{item.nome || item.nome_popular || "Sem nome"}</Text>
              <Text style={style.cardDescription} numberOfLines={3}>
                {item.descricao || item.especie || "Sem descrição disponível."}
              </Text>
          </View>
        </TouchableOpacity>
      </View>
    </View>
  );

  return (
    <SafeAreaView style={style.SafeAreaView}>
      <View style={style.busca}>
        <TextInput style={style.buscaInput} placeholder="Busca" />
      </View>

      {/* Conteúdo Dinâmico */}
      {loading ? (
        <View style={style.centerContainer}>
          <ActivityIndicator size="large" color="#007AFF" />
          <Text style={style.loadingText}>Carregando animais...</Text>
        </View>
      ) : error ? (
        <View style={style.centerContainer}>
          <Text style={style.errorText}>{error}</Text>
          <TouchableOpacity style={style.retryButton} onPress={carregarAnimais}>
            <Text style={style.retryButtonText}>Tentar Novamente</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <FlatList
          data={animais}
          keyExtractor={(item, index) => item.id ? item.id.toString() : index.toString()}
          renderItem={renderCardItem}
          contentContainerStyle={style.listContainer}
          showsVerticalScrollIndicator={false}
          ListEmptyComponent={
            <View style={style.centerContainer}>
              <Text style={style.emptyText}>Nenhum animal encontrado.</Text>
            </View>
          }
        />
      )}
    </SafeAreaView>
  );
};